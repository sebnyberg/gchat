package server

import (
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/sebnyberg/gchat/pkg/pb"
)

type server struct {
	Broadcast
}

func newChatServer() *server {
	return &server{
		Broadcast: NewBroadcast(),
	}
}

// ChatSession is single chat session between the server and the client.
// When a client sends a chat message to the server, it will be put in a common chat channel.
// This message will then be put in each clients own listen channel, which is sent back to the client
func (s *server) ChatSession(stream pb.ChatService_ChatSessionServer) error {
	log.Println("New client has connected")

	// Say hello to the client
	msg := &pb.ChatMessage{
		Username: "Server",
		Content:  "Welcome to the chat! Type a message and press enter to send a message.",
	}
	if err := stream.Send(&pb.ChatSessionResponse{Message: msg}); err != nil {
		fmt.Printf("Failed to send welcome message to client: %v\n", err)
	}

	// Create blocking channel
	waitc := make(chan struct{})

	chatc := s.Broadcast.GetChatc()

	uuid := fmt.Sprintf("%v", uuid.New())

	// Subscribe to messages from other clients
	listenc := s.Broadcast.Subscribe(uuid)
	defer s.Broadcast.Unsubscribe(uuid)

	// Receieve messages from client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("Client requested to disconnect from the server")
				break
			}
			if err != nil {
				log.Printf("Error when retrieving message from client: %v\n", err)
				break
			}
			chatc <- res.GetMessage()
		}

		fmt.Println("Done receiving messages from the client")
		close(waitc)
	}()

	// Send broadcasted messages to the client
	go func() {
		for msg := range listenc {
			err := stream.Send(&pb.ChatSessionResponse{Message: msg})
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}
	}()

	// Block until client disconnects
	<-waitc

	return nil
}

