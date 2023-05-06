package introspection_server

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/certainty/go-braces/internal/introspection/introspection_protocol"
)

type IntrospectionServer struct {
	EventChan    chan introspection_protocol.Event
	RequestChan  chan introspection_protocol.Request
	ResponseChan chan introspection_protocol.Response
	eventsSock   net.Listener
	controlSock  net.Listener
	tempDir      string
}

func NewServer() (*IntrospectionServer, error) {
	tempDir, err := ioutil.TempDir("", "introspection")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	eventSockPath := eventsSockPath(tempDir)
	controlSockPath := controlSockPath(tempDir)

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

	server := &IntrospectionServer{
		EventChan:    make(chan introspection_protocol.Event),
		RequestChan:  make(chan introspection_protocol.Request),
		ResponseChan: make(chan introspection_protocol.Response),
		eventsSock:   eventSock,
		controlSock:  controlSock,
		tempDir:      tempDir,
	}

	go server.handleEventConnections()
	go server.handleControlConnections()

	return server, nil
}

func (s *IntrospectionServer) handleEventConnections() {
	for {
		conn, err := s.eventsSock.Accept()
		if err != nil {
			continue
		}

		go s.handleEventStream(conn)
	}
}

func (s *IntrospectionServer) handleControlConnections() {
	for {
		conn, err := s.controlSock.Accept()
		if err != nil {
			continue
		}

		go s.handleControlRequests(conn)
		go s.handleControlResponses(conn)
	}
}

func (s *IntrospectionServer) handleEventStream(conn net.Conn) {
	enc := gob.NewEncoder(conn)
	for event := range s.EventChan {
		err := enc.Encode(introspection_protocol.WireEvent{Event: event})
		if err != nil {
			log.Printf("Error encoding event: %v", err)
		}
	}
}

func (s *IntrospectionServer) handleControlRequests(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	for {
		var req introspection_protocol.WireRequest
		err := dec.Decode(&req)
		if err != nil {
			log.Printf("Error decoding request: %v", err)
			continue
		} else {
			s.RequestChan <- req.Request
		}
	}
}

func (s *IntrospectionServer) handleControlResponses(conn net.Conn) {
	enc := gob.NewEncoder(conn)

	for response := range s.ResponseChan {
		err := enc.Encode(introspection_protocol.WireResponse{Response: response})

		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (s *IntrospectionServer) Close() {
	close(s.EventChan)
	close(s.RequestChan)
	close(s.ResponseChan)
	s.eventsSock.Close()
	s.controlSock.Close()
	os.Remove(controlSockPath(s.tempDir))
	os.Remove(eventsSockPath(s.tempDir))
	os.RemoveAll(s.tempDir)
}

func eventsSockPath(base string) string {
	return filepath.Join(base, "eventstream.ipc")
}

func controlSockPath(base string) string {
	return filepath.Join(base, "control.ipc")
}

func (s *IntrospectionServer) IPCDir() string {
	return s.tempDir
}
