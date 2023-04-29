package introspection

type IntrospectionEvent interface{}

type IntrospectionRequest interface{}

type IntrospectionResponse interface{}

type API interface {
	SendEvent(event IntrospectionEvent)
}

type IntrospectionChannel struct {
	events    chan IntrospectionEvent
	requests  chan IntrospectionRequest
	responses chan IntrospectionResponse
}

func (c IntrospectionChannel) SendEvent(event IntrospectionEvent) {
	c.events <- event
}

func NewChannel() IntrospectionChannel {
	return IntrospectionChannel{
		events:    make(chan IntrospectionEvent),
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
	return
}
