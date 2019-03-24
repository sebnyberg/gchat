package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sebnyberg/gchat/pkg/pb"
)

// Try to connect to the gchat server with the requested username
func joinServer(c pb.ChatServiceClient, username string) error {
	fmt.Printf("Joining chat as \"%v\"\n", username)

	// Create a context which times out after 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Connect with passed username
	connectionRequest := &pb.JoinServerRequest{Username: username}
	res, err := c.JoinServer(ctx, connectionRequest)
	if err != nil {
		return err
	}

	// Print response from server
	fmt.Printf("[Server]: %v", res.GetResponse())

	return nil
}

func chat(c pb.ChatServiceClient, username string) error {
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

		reader := bufio.NewReader(os.Stdin)
		for {

			text, err := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if err != nil {
				break
			}

			err = stream.Send(NewMessage(username, text))
			if err != nil {
				log.Fatalf("Failed to send message to the Server: %v", err)
				break
			}
		}

		stream.CloseSend()
		close(waitc)
	}()

	<-waitc

	return nil
}
