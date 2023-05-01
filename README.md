# Go Chat
An experiment to build a chat application using NATS support for websockets modeled after a pub-sub system

## How to run
Directly with Go using the following commnand:
```
go run cmd/main.go
```

### Docker
A docker file is included to easily create an image. To build the image issue the following command:
```
docker build -t gochat .
```

Once the image is built, run it with
```
docker run --rm --name gochat -p 3000:3000 -v $(pwd)/resources/config.yaml:/resources/config.yaml gochat
```
