package introspection

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	EventChan    chan IntrospectionEvent
	RequestChan  chan Request
	ResponseChan chan Response
	eventsSock   net.Listener
	controlSock  net.Listener
	tempDir      string
}

func NewServer() (*Server, error) {
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

	server := &Server{
		EventChan:    make(chan IntrospectionEvent),
		RequestChan:  make(chan Request),
		ResponseChan: make(chan Response),
		eventsSock:   eventSock,
		controlSock:  controlSock,
		tempDir:      tempDir,
	}

	go server.handleEventConnections()
	go server.handleControlConnections()

	return server, nil
}

func (s *Server) handleEventConnections() {
	for {
		conn, err := s.eventsSock.Accept()
		if err != nil {
			continue
		}

		go s.handleEventStream(conn)
	}
}

func (s *Server) handleControlConnections() {
	for {
		conn, err := s.controlSock.Accept()
		if err != nil {
			continue
		}

		go s.handleControlRequests(conn)
		go s.handleControlResponses(conn)
	}
}

func (s *Server) handleEventStream(conn net.Conn) {
	enc := gob.NewEncoder(conn)
	for event := range s.EventChan {
		err := enc.Encode(WireEvent{IntrospectionEvent: event})
		if err != nil {
			log.Printf("Error encoding event: %v", err)
		}
	}
}

func (s *Server) handleControlRequests(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	for {
		var req WireRequest
		err := dec.Decode(&req)
		if err != nil {
			log.Printf("Error decoding request: %v", err)
			continue
		} else {
			s.RequestChan <- req.Request
		}
	}
}

func (s *Server) handleControlResponses(conn net.Conn) {
	enc := gob.NewEncoder(conn)

	for response := range s.ResponseChan {
		err := enc.Encode(WireResponse{Response: response})

		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (s *Server) Close() {
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

func (s *Server) IPCDir() string {
	return s.tempDir
}
