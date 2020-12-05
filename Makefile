build:
	go build -o api cmd/api/main.go
	go build -o camera cmd/camera/main.go

api:
	go build -o api cmd/api/main.go

camera:
	go build -o camera cmd/camera/main.go

clean:
	rm api camera 

.PHONY: clean
