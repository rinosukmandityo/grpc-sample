package main

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/rinosukmandityo/grpc-sample/chat-sample/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var glog grpclog.LoggerV2

func init() {
	glog = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type connection struct {
	stream chat.BroadcastService_CreateStreamServer
	user   *chat.User
	active bool
	err    chan error
}

type server struct {
	connections []*connection
}

func (s *server) CreateStream(conn *chat.Connection, stream chat.BroadcastService_CreateStreamServer) error {
	c := &connection{
		stream: stream,
		user:   conn.GetUser(),
		active: conn.GetActive(),
		err:    make(chan error),
	}
	s.connections = append(s.connections, c)

	return <-c.err
}

func (s *server) BroadcastMessage(ctx context.Context, message *chat.Message) (*chat.Close, error) {
	wg := sync.WaitGroup{}
	done := make(chan int)
	for _, conn := range s.connections {
		wg.Add(1)

		go func(c *connection, msg *chat.Message) {
			defer wg.Done()

			if c.active {
				e := c.stream.Send(msg)
				glog.Infof("Sending a message to %s", c.user.GetName())
				if e != nil {
					c.active = false
					glog.Errorf("Failed to send a message to %s: %v", c.user.GetName(), e)
				}
			}
		}(conn, message)
	}

	go func() {
		wg.Wait()
		close(done)
	}()
	<-done

	return &chat.Close{}, nil
}

func main() {
	addr := ":8080"
	s := grpc.NewServer()
	srv := &server{[]*connection{}}
	chat.RegisterBroadcastServiceServer(s, srv)

	lis, e := net.Listen("tcp", addr)
	if e != nil {
		glog.Fatalf("could not connect to address: %v", e)
	}
	glog.Infof("Starting server at port %s", addr)
	glog.Fatal(s.Serve(lis))
}
