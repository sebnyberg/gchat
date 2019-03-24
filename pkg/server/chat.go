package server

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
var activeUsers map[string]bool = make(map[string]bool)

// Connect to the chat
// Returns AlreadyExists error if username is taken and that user is already connected
func (*server) JoinServer(ctx context.Context, req *pb.JoinServerRequest) (*pb.JoinServerResponse, error) {
	requestedUsername := req.GetUsername()
	log.Printf("Client requested to join server as \"%v\"", requestedUsername)

	// Check if the username is taken
	if _, ok := claimedUsernames[requestedUsername]; ok {

		// Check if claimed username is connected
		if _, ok := activeUsers[requestedUsername]; ok {
			// Ask client to re-try with a different name
			log.Printf("Username already taken: \"%v\"", requestedUsername)
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Failed to join server, username already taken: \"%v\"", requestedUsername))
		}

		// User is reconnecting
		response := &pb.JoinServerResponse{
			Response: fmt.Sprintf("Welcome back \"%v\"\n", requestedUsername),
		}

		return response, nil
	}

	// New client / username
	log.Printf("Adding username \"%v\" to list of claimed usernames\n", requestedUsername)
	claimedUsernames[requestedUsername] = true

	response := &pb.JoinServerResponse{
		Response: fmt.Sprintf("Successfully joined server as \"%v\"\n", requestedUsername),
	}

	return response, nil
}

func (*server) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("Chat initialized for client")

	waitc := make(chan struct{})

	// Receive messages
	go func() {
		var user string
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Printf("Client disconnected")
				break
			}
			if err != nil {
				log.Printf("Failed to receive message from client: %v", err)
				break
			}
			log.Printf("Broadcasting message to clients: %v\n", res.GetMessage())

			content := res.GetMessage().GetContent()
			user = res.GetMessage().GetUsername()

			if _, claimed := claimedUsernames[user]; !claimed {
				log.Println("Client not authenticated, please try again..")
				activeUsers[user] = false
				break
			}

			activeUsers[user] = true

			stream.Send(&pb.ChatResponse{
				Message: &pb.ChatMessage{
					Username: user,
					Content:  content,
				},
			})
		}
		if user != "" {
			claimedUsernames[user] = false
			activeUsers[user] = false
		}
		close(waitc)
	}()

	<-waitc

	return nil
}
