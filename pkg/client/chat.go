package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/sebnyberg/gchat/pkg/pb"
)

// Try to connect to the gchat server with the requested username
func connectToChat(c pb.ChatServiceClient, username string) error {
	fmt.Printf("Joining chat server with username %v\n", username)

	// Create a context which times out after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Connect with passed username
	connectionRequest := &pb.ConnectToChatRequest{Username: username}
	res, err := c.ConnectToChat(ctx, connectionRequest)
	if err != nil {
		return fmt.Errorf("Failed to connect to chat: %v", err)
	}

	// Print response from server
	fmt.Printf("[Server]: %v", res.GetResponse())

	return nil
}

func joinChat(c pb.ChatServiceClient) {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize chat stream: %v", err)
	}

	waitc := make(chan struct{})
	// Retrieve messages
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server closed the connection")
				break
			}
			if err != nil {
				log.Fatalf("Error when receiving message: %v", err)
				break
			}
			msg := res.GetMessage()
			log.Printf("%v: %v", msg.GetUsername(), msg.GetContent())
		}
		close(waitc)
	}()

	// Send messages
	go func() {
		msgs := []string{"Hello", "World", "Well this chat sucks... I'm out"}
		for _, msg := range msgs {
			err := stream.Send(NewMessage("Me", msg))
			if err != nil {
				log.Fatalf("Failed to send message to the Server: %v", err)
				break
			}
			time.Sleep(2 * time.Second)
		}
		close(waitc)
	}()

	<-waitc
}
