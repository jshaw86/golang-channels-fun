build:
	go build -o api cmd/api/main.go
	go build -o camera cmd/camera/main.go

clean:
	rm api camera 

.PHONY: clean
