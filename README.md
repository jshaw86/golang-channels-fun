# Summary
This program consists of a server (api) and a client (camera). Camera and API interact over http/1.1. The workflow is as follows:
1. Camera opens connection to the `/control-channel`, and long polls until either the client timeout is raised by the http transport or a call is made to `/logs` to fetch the data from the camera.
2. When `/logs` is called a message is sent via golang server's `ControlChannel` to close the long poll and return 200 OK back to the client, and waits for the logs from the `DataChannel`
3. If the client receives 200 OK it will immediately reconnect and send to the `/send-stats` endpoint the logs at that instant concurrently.
4. Once the server receives the call to `/send-stats` it will pipe each message of the payload body through the  golang server's `DataChannel`
5. Concurrently the `/logs` endpoint is ingesting from the `DataChannel` while the `/send-stats` endpoint is producing until `/send-stats` sends a `EOL` signal which tells the `/logs` it may respond with the collected logs  
6. By this time the client has reconnected to the `/control-channel` and waiting for another call to `/logs`

NOTE: client timeout is set to 1 minute currently and will gracefully disconnect and reconnect to `/control-channel` if there's no activity on `/logs`
 
# Compiling
Without docker, from root:
```
make
```

# Running
Without docker, from root:

1. Open two terminals, tmux or screen sessions
2. Run API in one terminal
```
./api
```
3. Run Camera from the other
```
./camera
```

# The Client

The client consists of 3 event loops:
1. MessageGeneration
2. RequestLoop
3. ConnectLoop

## Message Generation

Infinite for loop that generates a new "event" for the camera every 10 seconds. The event is passed out through a channel and into the RequestLoop for aggregation. It's more idiomatic in go to use message passing than semaphores for potentially "shared" memory

## RequestLoop

Infinte loop that listens for a new message from MessageGeneration or a 200 OK from ConnectLoop channel.

If it receives a new message from MessageGeneration it will append to a local slice. If it receives a message from ConnectLoop it will call `/send-stats`.

## ConnectLoop

Infinte for loop that connects to the `/control-channel` endpoint via client `Connect` method. If it receives a 200 OK from server will pass onto channel. If it receives a timeout it will reconnect.

# The Server

Consists of 3 endpoints:
1. `/logs`
2. `/control-channel`
3. `/send-stats`
and two internal channels:
1. `ControlChannel`
2. `DataChannel`

## Logs Endpoint

Logs signals to `/control-channel` via internal `ControlChannel` to close an existing open camera connection and waits for logs from the `DataChannel`

Once logs starts to be received on the `DataChannel` the logs endpoint will iterate until it sees an `EOL` at which it will write the response and close.

Returns JSON response that's an array of objects where each object has a `description` and `timestamp` field 

## Control Channel Endpoint

Waits for `send-stats` signal from `ControlChannel` or for client to close the connection

Returns 200 OK on success

## Send Stats Endpoint

Accepts payload of logs from camera, deserializes into `messages.Message` and passes onto `DataChannel` to be consumed by the the logs endpoint

Returns 200 OK on success


