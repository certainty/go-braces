package introspection

import "encoding/gob"

type MessageType uint8
type Payload interface{}

const (
	MessageTypeEvent MessageType = iota
	MessageTypeControl
	MessageTypeShutdown
)

type WireMessage struct {
	MessageType MessageType
	Payload     Payload
}

type Shutdown struct{}
type Hello struct{}

func RegisterWireTypes() {
	gob.Register(Shutdown{})
	gob.Register(Hello{})
	gob.Register(WireMessage{})
	gob.Register(MessageTypeEvent)
	gob.Register(MessageTypeControl)
	gob.Register(MessageTypeShutdown)
}
