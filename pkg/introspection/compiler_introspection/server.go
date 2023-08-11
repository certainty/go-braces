package compiler_introspection

import (
	"errors"
	"github.com/certainty/go-braces/pkg/introspection"
)

const INTROSPECTION_TOOL_NAME = "compiler-introspection-server"

type Server struct {
	wireServer *introspection.WireServer[CompilerIntrospectionControl, CompilerIntrospectionEvent]
	events     *introspection.WireEventConnection[CompilerIntrospectionEvent]
	control    *introspection.WireControlConnection[CompilerIntrospectionControl]
}

func NewServer() (*Server, error) {
	wireServer, err := introspection.NewWireServer[CompilerIntrospectionControl, CompilerIntrospectionEvent](INTROSPECTION_TOOL_NAME)
	if err != nil {
		return nil, err
	}
	return &Server{wireServer, nil, nil}, nil
}

func (s *Server) WaitForClient() (Instrumentation, error) {
	if s.events != nil || s.control != nil {
		return nil, errors.New("Server already has a client")
	}

	control, events, err := s.wireServer.Accept()
	if err != nil {
		return nil, err
	}

	s.events = events
	s.control = control

	return NewInstrumentationFromServer(s), nil
}

func (s *Server) Close() {
	if s.control != nil {
		s.control.Close()
	}

	if s.events != nil {
		s.events.Close()
	}
}

func (s *Server) IsConnected() bool {
	return s.events.IsOpen() && s.control.IsOpen()
}

// send events non blocking
func (s *Server) SendEvents(events ...CompilerIntrospectionEvent) error {
	if !s.IsConnected() {
		return errors.New("No client connected")
	}

	for _, event := range events {
		s.events.Channel <- event
	}

	return nil
}

func (s *Server) ReceiveControl() (CompilerIntrospectionControl, error) {
	if !s.IsConnected() {
		return nil, errors.New("No client connected")
	}
	control := <-s.control.In
	return control, nil
}

func (s *Server) SendControl(control CompilerIntrospectionControl) error {
	if !s.IsConnected() {
		return errors.New("No client connected")
	}
	s.control.Out <- control
	return nil
}
