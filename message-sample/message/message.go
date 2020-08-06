package message

import (
	"context"
	"log"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, msg *Text) (*Text, error) {
	log.Printf("Received message from client: %s", msg.GetBody())
	return &Text{Body: "Hello from the Server"}, nil
}
