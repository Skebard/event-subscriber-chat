package chat

import (
	"github.com/Skebard/event-subscriber-chat/src/pkg/pubsub"
)

type Message struct {
	Username string
	Content  string
	Time     int64
}

type RoomInterface interface {
	SendMessage(message Message)
	HandleIncomingMessages(HandleNewMessageCallback)
}

type HandleNewMessageCallback func(message Message)

type Room struct {
	Identifier      string
	Capacity        int
	CurrentCapacity int
	eventBus        pubsub.EventBusInterface
	Conversation    []Message
}

var rooms map[string]*Room

func NewRoom(identifier string, capacity int, eb pubsub.EventBusInterface) (*Room, error) {
	if _, ok := rooms[identifier]; ok {
		return nil, NewRoomIdentifierAlreadyInUse(identifier)
	}
	room := &Room{identifier, capacity, 1, eb, []Message{}}
	roomManager := newRoomManager(room)
	eb.Subscribe(room.Identifier, roomManager.ch, nil)
	if rooms == nil {
		rooms = map[string]*Room{}
	}
	rooms[identifier] = room
	return room, nil
}

func EnterRoom(roomIdentifier string, eb pubsub.EventBusInterface) (*Room, error) {
	if room, ok := rooms[roomIdentifier]; ok {
		if room.CurrentCapacity == room.Capacity {
			return nil, NewExceedRoomCapacityError(roomIdentifier, room.Capacity)
		}
		room.CurrentCapacity++
		return room, nil
	}
	return nil, NewRoomNotFoundError(roomIdentifier)
}

func LeaveRoom(room *Room) error {
	room.CurrentCapacity--
	return nil
}

func (room *Room) SendMessage(message Message) {
	room.eventBus.Publish(pubsub.NewDataEvent(room.Identifier, message), nil)
}

func (room *Room) HandleIncomingMessages(handle HandleNewMessageCallback) error {
	ch := make(chan pubsub.DataEvent)
	room.eventBus.Subscribe(room.Identifier, ch, nil)
	go func() {
		for {
			event := <-ch
			if msg, ok := event.Data.(Message); ok {
				handle(msg)
			}
		}
	}()
	return nil
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
