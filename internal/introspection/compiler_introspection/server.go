package compiler_introspection

import (
	"errors"
	"github.com/certainty/go-braces/internal/introspection"
)

const INTROSPECTION_TOOL_NAME = "compiler-introspection-server"

type Server struct {
	wireServer *introspection.WireServer
	events     *introspection.WireEventConnection
	control    *introspection.WireControlConnection
}

func NewServer() (*Server, error) {
	wireServer, err := introspection.NewWireServer(INTROSPECTION_TOOL_NAME)
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

func (s *Server) SendEventsSync(events ...CompilerIntrospectionEvent) error {
	if !s.IsConnected() {
		return errors.New("No client connected")
	}

	for _, event := range events {
		s.events.Channel <- event
	}

	return nil
}
