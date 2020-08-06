package main

import (
	"log"
	"net"

	"github.com/rinosukmandityo/grpc-sample/message-sample/message"

	"google.golang.org/grpc"
)

func main() {
	addr := ":8989"
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, e)
	}
	s := grpc.NewServer()
	srv := message.Server{}
	message.RegisterMessageServiceServer(s, &srv)

	log.Println("Listening to", addr)
	log.Fatal(s.Serve(lis))
}
