package client

import (
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
	"google.golang.org/grpc"
)

func NewMessage(author string, content string) *pb.ChatRequest {
	return &pb.ChatRequest{
		Message: &pb.ChatMessage{
			Author:  author,
			Content: content,
		},
	}
}

func ConnectClient() error {
	log.Println("Creating a new client...")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer cc.Close()

	c := pb.NewChatServiceClient(cc)

	joinChat(c)

	return nil
}
