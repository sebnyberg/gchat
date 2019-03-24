package server

import (
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
)

type server struct{}

func NewChatServer() *server {
	return &server{}
}

var connectedResponse *pb.ChatResponse = &pb.ChatResponse{
	Message: &pb.ChatMessage{
		Author:  "TODO name",
		Content: "Connected to chat as TODO name",
	},
}

func (*server) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("Client connected as HEHE XD")
	err := stream.Send(connectedResponse)
	if err != nil {
		log.Fatalf("Failed to send connected response: %v", err)
	}

	waitc := make(chan struct{})

	// Receive messages
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message from client: %v", err)
			close(waitc)
		}
		log.Printf("Received message from client: %v\n", msg.GetMessage())
	}

	<-waitc

	return nil
}
