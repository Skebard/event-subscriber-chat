package chat

import (
	"testing"

	"github.com/Skebard/event-subscriber-chat/src/pkg/pubsub"
)

const (
	VALID_TEST_IDENTIFIER string = "identifier-1"
	VALID_TEST_CAPACITY   int    = 5
	VALID_TEST_USERNAME   string = "test-username"
)

func TestNewRoomValidParams(t *testing.T) {
	eventBus := pubsub.NewEventBus()
	if _, err := NewRoom(VALID_TEST_IDENTIFIER, VALID_TEST_CAPACITY, eventBus); err != nil {
		t.Fatalf("Not possible to create room")
	}
}

func TestRoomIdentifierAlreadyExists(t *testing.T) {
	eventBus := pubsub.NewEventBus()
	NewRoom(VALID_TEST_IDENTIFIER, VALID_TEST_CAPACITY, eventBus)
	if _, err := NewRoom(VALID_TEST_IDENTIFIER, VALID_TEST_CAPACITY, eventBus); err != nil {
		if _, ok := err.(*RoomIdentifierAlreadyInUse); !ok {
			t.Fatalf("Invalid error returned")
		}
	} else {
		t.Fatalf("Room created")
	}
}

// func TestRoomSendMessage(t *testing.T) {
// 	eventBus := pubsub.NewEventBus()
// 	room, _ := NewRoom(VALID_TEST_IDENTIFIER, VALID_TEST_CAPACITY, eventBus)
// 	for i := 0; i < 10; i++ {
// 		msg := fmt.Sprintf("Message %d", i)
// 		room.SendMessage(Message{VALID_TEST_USERNAME, msg + "2", time.Now().UnixMilli()})
// 		if room.Conversation[i].Content != msg {
// 			t.Fatalf("Message %s not received", msg)
// 		}
// 	}
// }
