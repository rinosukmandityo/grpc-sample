FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash git

RUN mkdir /chat-sample
RUN mkdir -p /chat-sample/chat 

WORKDIR /chat-sample

COPY ./chat/chat.pb.go ./chat
COPY ./cmd/server/server.go .

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o chat-sample server.go

CMD ./chat-sample