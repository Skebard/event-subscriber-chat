package chat

import "github.com/Skebard/event-subscriber-chat/src/pkg/pubsub"

type Message struct {
	Username string
	Content  string
	Time     int
}

type RoomInterface interface {
	SendMessage(message Message)
	HandleIncomingMessages(HandleNewMessageCallback)
}

type HandleNewMessageCallback func(message Message)

type Room struct {
	Identifier   string
	Capacity     int
	eventBus     pubsub.EventBusInterface
	Conversation []Message
}

func NewRoom(identifier string, capacity int, eb pubsub.EventBusInterface) *Room {
	room := &Room{identifier, capacity, eb, []Message{}}
	roomManager := newRoomManager(room)
	eb.Subscribe(room.Identifier, roomManager.ch, nil)
	return room
}

// func EnterRoom(roomIdentifier string, eb pubsub.EventBusInterface) (*Room, error) {
// 	eb.

// }

func (room *Room) SendMessage(message Message) {
	room.eventBus.Publish(pubsub.NewDataEvent(room.Identifier, message), nil)
}

type roomManager struct {
	ch   pubsub.DataChannel
	room *Room
}

func newRoomManager(room *Room) *roomManager {
	ch := make(chan pubsub.DataEvent)
	roomManager := &roomManager{ch, room}
	go roomManager.lookForEvents()
	return roomManager
}

func (roomManager *roomManager) lookForEvents() {
	for {
		event := <-roomManager.ch
		roomManager.handleNewEvent(event)
	}
}

func (roomManager *roomManager) handleNewEvent(event pubsub.DataEvent) {
	if msg, ok := event.Data.(Message); ok {
		roomManager.handleNewMessage(msg)
	}
}

func (roomManager *roomManager) handleNewMessage(message Message) {
	roomManager.addMessageToConversation(message, roomManager.room)
}

func (roomManager *roomManager) addMessageToConversation(message Message, room *Room) {
	room.Conversation = append(room.Conversation, message)
}
