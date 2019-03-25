package client

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sebnyberg/gchat/pkg/pb"
)

func chat(c pb.ChatServiceClient, username string) error {
	stream, err := c.ChatSession(context.Background())

	if err != nil {
		log.Fatalf("Failed to initialize chat stream: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		listenToChat(stream)

		close(waitc)
		stream.CloseSend()
	}()

	go func() {
		sendMessages(stream, username)

		close(waitc)
		stream.CloseSend()
	}()

	<-waitc

	return nil
}

func listenToChat(stream pb.ChatService_ChatSessionClient) {
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
}

func sendMessages(stream pb.ChatService_ChatSessionClient, username string) {
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
}

