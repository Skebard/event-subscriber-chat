package pubsub

import (
	"sync"
)

type EventBusInterface interface {
	Publish(event DataEvent, ctx interface{}) error
	Subscribe(topic string, ch DataChannel, ctx interface{}) error
	CreatePrivateTopic(topic string, ctx interface{}, authFunc AuthenticatePrivateChannel)
}

func NewEventBus() EventBusInterface {
	return &EventBus{
		subscribers:   map[string]DataChannelSlice{},
		privateTopics: map[string]PrivateTopic{},
	}
}

type EventBus struct {
	subscribers   map[string]DataChannelSlice
	privateTopics map[string]PrivateTopic
	rm            sync.RWMutex
}

func (eb *EventBus) Publish(event DataEvent, ctx interface{}) error {
	eb.rm.RLock()
	defer eb.rm.RUnlock()
	if chans, found := eb.subscribers[event.Topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		// special thanks for /u/freesid who pointed it out
		channels := append(DataChannelSlice{}, chans...)
		go eb.dispatchEvent(event, channels)
	} else if privateTopic, found := eb.privateTopics[event.Topic]; found {
		privateTopic.Authenticate(event, eb.privateTopics[event.Topic].Ctx, ctx)
		channels := append(DataChannelSlice{}, privateTopic.Channels...)
		go eb.dispatchEvent(event, channels)
	}
	return nil
}

func (eb *EventBus) dispatchEvent(data DataEvent, dataChannels DataChannelSlice) {
	for _, ch := range dataChannels {
		ch <- data
	}
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel, ctx interface{}) error {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else if privateTopic, found := eb.privateTopics[topic]; found {
		if !privateTopic.Authenticate(DataEvent{}, privateTopic.Ctx, ctx) {
			return nil
		}
		privateTopic.Channels = append(privateTopic.Channels, ch)
		eb.privateTopics[topic] = privateTopic
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	return nil
}

func (eb *EventBus) CreatePrivateTopic(topic string, ctx interface{}, authFunc AuthenticatePrivateChannel) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()
	if _, found := eb.subscribers[topic]; found {
		return
	}
	if _, found := eb.privateTopics[topic]; found {
		return
	}
	eb.privateTopics[topic] = PrivateTopic{
		Channels:     DataChannelSlice{},
		Authenticate: authFunc,
		Name:         topic,
		Ctx:          ctx,
	}
}
