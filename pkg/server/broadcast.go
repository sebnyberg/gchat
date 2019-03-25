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

<<<<<<< HEAD
	// Initiate broadcast
=======
>>>>>>> 27157ccf12aeb8e4ad5b2dd09146e5aa985ceca0
	b.Start()

	return b
}

<<<<<<< HEAD
// GetChatc retrives a channel which is used to gather messages from clients
=======
// The Chatc is a common channel that stores all messages from the clients
>>>>>>> 27157ccf12aeb8e4ad5b2dd09146e5aa985ceca0
func (b *Broadcast) GetChatc() chan *pb.ChatMessage {
	return b.Chatc
}

// Subscribe registers a subscriber and returns a channel which can be
// listened to for broadcasted messages
func (b *Broadcast) Subscribe(subscriber string) chan *pb.ChatMessage {
<<<<<<< HEAD
	log.Println("Subscribing client to messages from other clients")

	// Make sure that broadcasting is enabled
	if !b.IsBroadcasting {
		log.Println("Currently not broadcasting messages, did you run broadcast.Start()?")
=======
	log.Println("Adding client to list of subscribers")
	// Make sure that broadcasting has started
	if !b.IsBroadcasting {
		log.Println("Currently not broadcasting, did you run broadcast.Start()?")
>>>>>>> 27157ccf12aeb8e4ad5b2dd09146e5aa985ceca0
	}

	s := make(chan *pb.ChatMessage)

	b.Subscribers[subscriber] = s

	return s
}

// Unsubscribe removes a subscriber from the broadcast
func (b *Broadcast) Unsubscribe(subscriber string) {
<<<<<<< HEAD
	log.Println("Unsubscribing client from broadcast")
=======
	log.Println("Removing subscriber from the broadcast")
>>>>>>> 27157ccf12aeb8e4ad5b2dd09146e5aa985ceca0
	// Check if subscriber exists
	if _, ok := b.Subscribers[subscriber]; !ok {
		log.Printf("Tried to unsubscribe a subscriber that did not exist: %v", subscriber)
		return
	}

	delete(b.Subscribers, subscriber)
}

<<<<<<< HEAD
// Start broadcasting messages from the chat channel (chatc) to all subscriber channels
=======
// Start broadcasting all messages that appear in the chatc subscriber channels
>>>>>>> 27157ccf12aeb8e4ad5b2dd09146e5aa985ceca0
func (b *Broadcast) Start() {
	go func() {
		for msg := range b.Chatc {
			for _, sub := range b.Subscribers {
				sub <- msg
			}
		}
	}()
}

