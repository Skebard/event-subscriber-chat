package pubsub

type PrivateTopic struct {
	Channels     DataChannelSlice
	Authenticate AuthenticatePrivateChannel
	Name         string
	Ctx          interface{}
}

type AuthenticatePrivateChannel func(event DataEvent, topicCtx interface{}, ctx interface{}) bool
