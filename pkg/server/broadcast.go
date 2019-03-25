package server

import (
	"log"

	"github.com/sebnyberg/gchat/pkg/pb"
)

type Broadcast struct {
	Chatc          chan *pb.ChatMessage
	Subscribers    map[string]chan *pb.ChatMessage
	IsBroadcasting bool
}

func NewBroadcast() Broadcast {
	b := Broadcast{
		Chatc:          make(chan *pb.ChatMessage),
		Subscribers:    make(map[string]chan *pb.ChatMessage),
		IsBroadcasting: true,
	}

	b.Start()

	return b
}

// The Chatc is a common channel that stores all messages from the clients
func (b *Broadcast) GetChatc() chan *pb.ChatMessage {
	return b.Chatc
}

// Subscribe registers a subscriber and returns a channel which can be
// listened to for broadcasted messages
func (b *Broadcast) Subscribe(subscriber string) chan *pb.ChatMessage {
	log.Println("Adding client to list of subscribers")
	// Make sure that broadcasting has started
	if !b.IsBroadcasting {
		log.Println("Currently not broadcasting, did you run broadcast.Start()?")
	}

	s := make(chan *pb.ChatMessage)

	b.Subscribers[subscriber] = s

	return s
}

// Unsubscribe removes a subscriber from the broadcast
func (b *Broadcast) Unsubscribe(subscriber string) {
	log.Println("Removing subscriber from the broadcast")
	// Check if subscriber exists
	if _, ok := b.Subscribers[subscriber]; !ok {
		log.Printf("Tried to unsubscribe a subscriber that did not exist: %v", subscriber)
		return
	}

	delete(b.Subscribers, subscriber)
}

// Start broadcasting all messages that appear in the chatc subscriber channels
func (b *Broadcast) Start() {
	go func() {
		for msg := range b.Chatc {
			for _, sub := range b.Subscribers {
				sub <- msg
			}
		}
	}()
}

