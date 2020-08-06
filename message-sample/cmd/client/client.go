package main

import (
	"context"
	"log"

	"github.com/rinosukmandityo/grpc-sample/message-sample/message"

	"google.golang.org/grpc"
)

func main() {
	addr := ":8989"
	cc, e := grpc.Dial(addr, grpc.WithInsecure())
	if e != nil {
		log.Fatalf("could not connect to %s: %v", addr, e)
	}
	defer cc.Close()

	client := message.NewMessageServiceClient(cc)
	msg, e := client.SayHello(context.Background(), &message.Text{Body: "Message from the Client"})
	if e != nil {
		log.Fatalf("could not call SayHello procedure: %v", e)
	}
	log.Printf("Received message from server: %s", msg.GetBody())
}
