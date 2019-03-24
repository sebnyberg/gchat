package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/sebnyberg/gchat/pkg/pb"
)

func joinChat(c pb.ChatServiceClient) {
	log.Println("Connecting to server...")

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
			log.Printf("%v: %v", msg.GetAuthor(), msg.GetContent())
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
