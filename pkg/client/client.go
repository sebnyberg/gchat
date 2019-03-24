package client

import (
	"fmt"

	"github.com/sebnyberg/gchat/pkg/pb"
	"google.golang.org/grpc"
)

func NewMessage(username string, content string) *pb.ChatRequest {
	return &pb.ChatRequest{
		Message: &pb.ChatMessage{
			Username: username,
			Content:  content,
		},
	}
}

func RunClient() error {
	fmt.Println("Creating a new client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	c := pb.NewChatServiceClient(cc)

	// Try to join chat as "anonymous"
	username := "anonymous"
	err = connectToChat(c, username)
	if err != nil {
		return fmt.Errorf("Failed to connect to the server: %v", err)
	}

	return nil
}
