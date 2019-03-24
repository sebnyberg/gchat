package server

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
)

type server struct{}

func NewChatServer() *server {
	return &server{}
}

var connectedResponse *pb.ChatResponse = &pb.ChatResponse{
	Message: &pb.ChatMessage{
		Username: "TODO name",
		Content:  "Connected to chat as TODO name",
	},
}

var claimedUsernames map[string]bool = make(map[string]bool)
var connectedUsers map[string]bool = make(map[string]bool)

func (*server) ConnectToChat(ctx context.Context, req *pb.ConnectToChatRequest) (*pb.ConnectToChatResponse, error) {
	requestedUsername := req.GetUsername()
	log.Printf("Client requested to join chat with username %v", requestedUsername)

	// Check if the username is taken
	if _, ok := claimedUsernames[requestedUsername]; ok {

		// Check if claimed username is connected
		if _, ok := connectedUsers[requestedUsername]; ok {
			// Ask client to re-try with a different name
			log.Printf("Username already taken: %v", requestedUsername)
			return nil, fmt.Errorf("Failed to connect, username already taken: %v", requestedUsername)
		}

		// User is reconnecting
		response := &pb.ConnectToChatResponse{
			Response: fmt.Sprintf("Welcome back %v\n", requestedUsername),
		}

		return response, nil
	}

	// New client / username
	log.Printf("Adding username %v to the list of claimed usernames\n", requestedUsername)
	claimedUsernames[requestedUsername] = true

	response := &pb.ConnectToChatResponse{
		Response: fmt.Sprintf("Successfully claimed username %v\n", requestedUsername),
	}

	return response, nil
}

func (*server) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("Chat initialized for client")

	waitc := make(chan struct{})

	// Receive messages
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Client disconnected")
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive message from client: %v", err)
				break
			}
			log.Printf("Received message from client: %v\n", msg.GetMessage())
		}
		close(waitc)
	}()

	<-waitc

	return nil
}
