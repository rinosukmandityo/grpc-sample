package main

import (
	"context"
	"log"
	"net"

	calc "github.com/rinosukmandityo/grpc-sample/calculation-sample/calculation"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Add(ctx context.Context, req *calc.Request) (*calc.Response, error) {
	return &calc.Response{Result: req.GetA() + req.GetB()}, nil
}

func (s *server) Multiply(ctx context.Context, req *calc.Request) (*calc.Response, error) {
	return &calc.Response{Result: req.GetA() * req.GetB()}, nil
}

func main() {
	addr := ":8989"
	s := grpc.NewServer()

	calc.RegisterCalculationServiceServer(s, &server{})
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		log.Println("could not listen to address %s: %v", addr, e)
	}

	log.Println("Listening to", addr)
	log.Fatal(s.Serve(lis))
}
