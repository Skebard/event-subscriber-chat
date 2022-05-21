package pubsub

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent

type DataChannelSlice []DataChannel

func NewDataEvent(topic string, data interface{}) DataEvent {
	return DataEvent{data, topic}
}
