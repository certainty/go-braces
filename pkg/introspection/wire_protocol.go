package introspection

type MessageType uint8

const (
	MessageTypeEvent MessageType = iota
	MessageTypeControl
	MessageTypeShutdown
)

type WireMessage[T any] struct {
	MessageType MessageType
	Payload     T
}

type Shutdown struct{}
type Hello struct{}
