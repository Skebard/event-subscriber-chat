package chat

import "github.com/Skebard/event-subscriber-chat/src/pkg/pubsub"

type PublicRoom struct {
	*Room
}

func NewPublicRoom(name string, capacity int, eb pubsub.EventBusInterface) (*PublicRoom, error) {
	if room, error := NewRoom(name, capacity, eb); error != nil {
		return nil, error
	} else {
		return &PublicRoom{room}, nil
	}
}
