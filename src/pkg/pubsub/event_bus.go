package pubsub

import "sync"

type EventBusInterface interface {
	Publish(event DataEvent) error
	Subscribe(topic string, ch DataChannel) error
}

func NewEventBus() EventBusInterface {
	return &EventBus{
		subscribers: map[string]DataChannelSlice{},
	}
}

type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func (eb *EventBus) Publish(event DataEvent) error {
	eb.rm.RLock()
	defer eb.rm.RUnlock()
	if chans, found := eb.subscribers[event.Topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		// special thanks for /u/freesid who pointed it out
		channels := append(DataChannelSlice{}, chans...)
		go eb.dispatchEvent(event, channels)
	}
	return nil
}

func (eb *EventBus) dispatchEvent(data DataEvent, dataChannels DataChannelSlice) {
	for _, ch := range dataChannels {
		ch <- data
	}
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) error {
	eb.rm.Lock()
	defer eb.rm.Unlock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
		return nil
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
		return nil
	}
}
