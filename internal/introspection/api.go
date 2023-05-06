package introspection

import "fmt"

type IntrospectionEvent interface {
	EventInspect() string
}

type BeginCompileStringEvent struct {
	Input string
}

func (e BeginCompileStringEvent) EventInspect() string {
	return fmt.Sprintf("(BeginCompileStringEvent %s)", e.Input)
}

type IntrospectionRequest interface{}

type IntrospectionResponse interface{}

type API interface {
	SendEvent(event IntrospectionEvent)
	ReceiveEvent() IntrospectionEvent
}

type IntrospectionChannel struct {
	events    chan IntrospectionEvent
	requests  chan IntrospectionRequest
	responses chan IntrospectionResponse
}

func (c IntrospectionChannel) SendEvent(event IntrospectionEvent) {
	c.events <- event
}

func (c IntrospectionChannel) ReceiveEvent() IntrospectionEvent {
	return <-c.events
}

func NewAPI() IntrospectionChannel {
	return IntrospectionChannel{
		events:    make(chan IntrospectionEvent, 512),
		requests:  make(chan IntrospectionRequest),
		responses: make(chan IntrospectionResponse),
	}
}

type Null struct{}

func NullAPI() Null {
	return Null{}
}

// implements API
func (n Null) SendEvent(event IntrospectionEvent) {
}

func (n Null) ReceiveEvent() IntrospectionEvent {
	return nil
}
