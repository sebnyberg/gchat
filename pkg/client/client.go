package client

import (
	"fmt"
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewMessage(username string, content string) *pb.ChatSessionRequest {
	return &pb.ChatSessionRequest{
		Message: &pb.ChatMessage{
			Username: username,
			Content:  content,
		},
	}
}

func StartClient(username string) {
	fmt.Println("Running client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server")
	}
	defer cc.Close()

	c := pb.NewChatServiceClient(cc)

	// Join the chat
	if err := chat(c, username); err != nil {
		handleChatError(err)
	}
}

func handleChatError(err error) {
	if statusErr, ok := status.FromError(err); ok {
		if statusErr.Code() == codes.AlreadyExists {
			fmt.Println("Failed to connect to the server, username is taken")
		} else {
			panic(fmt.Sprintf("Unexpected RPC error: %v", statusErr))
		}
	} else {
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}
}
