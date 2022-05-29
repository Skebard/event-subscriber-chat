package chat

import "github.com/Skebard/event-subscriber-chat/src/pkg/pubsub"

type PublicRoom struct {
	*Room
}

func NewPublicRoom(name string, capacity int, eb pubsub.EventBusInterface) *PublicRoom {
	return &PublicRoom{NewRoom(name, capacity, eb)}
}
