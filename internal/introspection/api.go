package introspection

import (
	"encoding/gob"
	"github.com/certainty/go-braces/internal/introspection/introspection_protocol"
	"github.com/certainty/go-braces/internal/introspection/introspection_server"
)

func RegisterTypes() {
	gob.Register(introspection_protocol.HeloRequest{})
	gob.Register(introspection_protocol.HeloResponse{})
	gob.Register(BeginCompileStringEvent{})
	gob.Register(StartCompileModuleEvent{})
	gob.Register(EndCompileModuleEvent{})
}

type IntrospectionEvent interface {
	EventInspect() string
}

type IntrospectionRequest interface{}

type IntrospectionResponse interface{}

// This is a sequential interface to the introspection capabiltities
// The code is instrumented with this API
type API interface {
	SendEvent(event introspection_protocol.Event)
	// Future functions
	// WaitSingleStep(state CurrentState) waits until the client resumes
	// Abort() askes the interface if execution should be aborted
}

type APIFromServer struct {
	sever *introspection_server.IntrospectionServer
}

func (c APIFromServer) SendEvent(event introspection_protocol.Event) {
	c.sever.EventChan <- event
}

func NewAPI(server *introspection_server.IntrospectionServer) APIFromServer {
	return APIFromServer{server}
}

type Null struct{}

func NullAPI() Null {
	return Null{}
}

// implements API
func (n Null) SendEvent(event introspection_protocol.Event) {
	return
}
