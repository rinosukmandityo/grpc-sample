package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/rinosukmandityo/grpc-sample/chat-sample/chat"

	"google.golang.org/grpc"
)

var (
	client chat.BroadcastServiceClient
	wg     *sync.WaitGroup
)

func init() {
	wg = new(sync.WaitGroup)
}

func connect(user *chat.User) error {
	var err error

	streamClient, e := client.CreateStream(context.Background(), &chat.Connection{
		User:   user,
		Active: true,
	})

	if e != nil {
		return fmt.Errorf("connection failed: %v", e)
	}

	wg.Add(1)

	go func(stream chat.BroadcastService_CreateStreamClient) {
		defer wg.Done()

		for {
			msg, e := stream.Recv()
			if e != nil {
				err = fmt.Errorf("Error reading message: %v", e)
				break
			}
			fmt.Printf("[%s] %s: %s\n", time.Unix(0, msg.GetTimestamp()).Format("15:04"), msg.GetUser().GetName(), msg.GetContent())
		}
	}(streamClient)

	return err
}

func main() {
	name := ""
	flag.StringVar(&name, "N", "", "The name of the user")
	address := ""
	flag.StringVar(&address, "addr", "", "The name of the user")
	flag.Parse()
	if name == "" {
		log.Fatal("Please provide a name after flag -N")
	}
	addr := ":8080"
	if address != "" {
		addr = address + addr
	}

	done := make(chan int)
	cc, e := grpc.Dial(addr, grpc.WithInsecure())
	if e != nil {
		log.Printf("Could not connect to server %s: %v", addr, e)
	}
	client = chat.NewBroadcastServiceClient(cc)

	timestamp := time.Now()
	id := sha256.Sum256([]byte(fmt.Sprintf("%s%s", timestamp.String(), name)))
	user := &chat.User{
		Id:   hex.EncodeToString(id[:]),
		Name: name,
	}

	if e = connect(user); e != nil {
		log.Fatal(e.Error())
	}
	log.Println("let's start chatting . . .")

	wg.Add(1)

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chat.Message{
				User:      user,
				Content:   scanner.Text(),
				Timestamp: time.Now().UnixNano(),
			}

			_, e := client.BroadcastMessage(context.Background(), msg)
			if e != nil {
				log.Printf("Error sending message %v", e)
				break
			}
		}
	}()

	go func() {
		wg.Wait()
		close(done)
	}()
	<-done
}
