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

// Chat is a single open connection between the server and the client
// All chat sessions share a common chat channel which is posted to whenever a new message
// arrives from a client to the server.
// Each chat session has its own listen channel which is used to broadcast messages to each client.
func (s *server) ChatSession(stream pb.ChatService_ChatSessionServer) error {
	log.Println("Chat initialized for client")

	waitc := make(chan struct{})

	chatc := s.Broadcast.GetChatc()

	// Create a unique id and subscribe with it
	uuid := fmt.Sprintf("%v", uuid.New())

	listenc := s.Broadcast.Subscribe(uuid)
	defer s.Broadcast.Unsubscribe(uuid)

	// Receieve messages from client
	go func() {
		for {
			res, err := stream.Recv()
			log.Printf("Receieved message: %v", res)
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

	// Broadcast messages sent to the server to the client
	for msg := range listenc {
		go func(m *pb.ChatMessage) {
			err := stream.Send(&pb.ChatSessionResponse{Message: m})
			if err != nil {
				log.Printf("Failed to send message to client: %v", err)
			}
		}(msg)
	}

	<-waitc

	return nil
}

