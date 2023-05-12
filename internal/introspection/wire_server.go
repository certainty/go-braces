package introspection

import (
	"fmt"
	"log"
	"net"
	"os"
)

type WireServer struct {
	scope string

	eventSocketPath string
	eventSock       net.Listener

	controlSock       net.Listener
	controlSocketPath string
}

func NewWireServer(scope string) (*WireServer, error) {
	SetupDirectories(scope)
	eventSockPath := EventSocketPath(scope)
	controlSockPath := ControlSocketPath(scope)

	os.Remove(eventSockPath)
	os.Remove(controlSockPath)

	eventSock, err := net.Listen("unix", eventSockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on event stream socket: %w", err)
	}

	controlSock, err := net.Listen("unix", controlSockPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on control socket: %w", err)
	}

	return &WireServer{
		eventSocketPath:   eventSockPath,
		controlSocketPath: controlSockPath,
		eventSock:         eventSock,
		controlSock:       controlSock,
		scope:             scope,
	}, nil
}

func (s *WireServer) Accept() (*WireControlConnection, *WireEventConnection, error) {
	for {
		controlCon, err := s.controlSock.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			return nil, nil, err
		}
		wireControlConnection := NewWireControlConnection(controlCon)

		eventsConn, err := s.eventSock.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			wireControlConnection.Close()
			return nil, nil, err
		}
		wireEventConnection := NewWireEventConnection(eventsConn, WireEventSource)

		return wireControlConnection, wireEventConnection, nil
	}
}
