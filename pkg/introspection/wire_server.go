package introspection

import (
	"fmt"
	"log"
	"net"
)

type WireServer[C any, E any] struct {
	scope string

	eventSocketPath string
	eventSock       net.Listener

	controlSock       net.Listener
	controlSocketPath string
}

func NewWireServer[C any, E any](scope string) (*WireServer[C, E], error) {
	if err := SetupDirectories(scope); err != nil {
		return nil, fmt.Errorf("failed to setup directories: %w", err)
	}

	eventSockPath := EventSocketPath(scope)
	controlSockPath := ControlSocketPath(scope)

	eventSock, err := net.Listen("unix", eventSockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on event stream socket: %w", err)
	}

	controlSock, err := net.Listen("unix", controlSockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on control socket: %w", err)
	}

	return &WireServer[C, E]{
		eventSocketPath:   eventSockPath,
		controlSocketPath: controlSockPath,
		eventSock:         eventSock,
		controlSock:       controlSock,
		scope:             scope,
	}, nil
}

func (s *WireServer[C, E]) Accept() (*WireControlConnection[C], *WireEventConnection[E], error) {
	controlCon, err := s.controlSock.Accept()
	if err != nil {
		log.Printf("Error accepting connection: %v", err)
		return nil, nil, err
	}
	wireControlConnection := NewWireControlConnection[C](controlCon)

	eventsConn, err := s.eventSock.Accept()
	if err != nil {
		log.Printf("Error accepting connection: %v", err)
		wireControlConnection.Close()
		return nil, nil, err
	}
	wireEventConnection := NewWireEventConnection[E](eventsConn, WireEventSource)

	return wireControlConnection, wireEventConnection, nil
}
