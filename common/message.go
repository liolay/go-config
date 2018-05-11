package common

const (
	ClientConnect      MessageType = 1
	ClientConnectReply MessageType = 2
	ServerPushFile     MessageType = 3
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

func NewServerPushFileMessage(data []byte) *Message {
	return &Message{ServerPushFile, data}
}

type ServerPushedFile struct {
	App string
	Name string
	Content []byte
}

type MessageType int
