# Chat Sample
Sample of broadcast chat message application using gRPC

How to run
---

```bash
# clone this repo
git clone https://github.com/rinosukmandityo/grpc-sample.git

# go to the project directory
cd grpc-sample/chat-sample
```

## Run the server via CLI

```bash
# go to the server directory
cd cmd/server

# start the server app
go run server.go
```

## Run the server via Docker

```bash
# build the docker image
docker build -t chat-sample -f dockerfile.Dockerfile .

# run the docker container
docker run -it -p 8080:8080 chat-sample
```

## Run the client

```bash
# go to the client directory
cd cmd/client

# start the client app
go run client.go -N <user-name>

# start the client app with Docker, add <docker_ip> into flag -addr, e.g: 192.168.99.100
go run client.go -N <user-name> -addr <docker_ip>
```