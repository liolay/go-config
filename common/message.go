package common

const (
	ClientConnect MessageType = 1
	ClientConnectReply MessageType = 2
)

type Message struct {
	MessageType MessageType
	Data        []byte
}

func NewClientConnectMessage(data []byte) *Message {
	return &Message{ClientConnect, data}
}

func NewClientConnectReplyMessage(data []byte) *Message {
	return &Message{ClientConnectReply, data}
}

type MessageType int
