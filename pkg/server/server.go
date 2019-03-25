package server

import (
	"log"
	"net"

	"github.com/sebnyberg/gchat/pkg/pb"
	"google.golang.org/grpc"
)

func StartServer() error {
	log.Println("Starting server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	srv := newChatServer()
	pb.RegisterChatServiceServer(s, srv)

	log.Println("Chat server up and listening for connections")
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}
