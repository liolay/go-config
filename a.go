package main


const (
	ClientConnect = 0
	ServerSendFile = 0
)

type Message interface {
	Type() int
}

type ConnectMessage struct {
}

func (c *ConnectMessage) Type() int {
	return ClientConnect
}

func NewConnectMessage() *ConnectMessage  {
	return &ConnectMessage{}
}
