package introspection_server

import (
	"errors"
	"github.com/certainty/go-braces/internal/introspection/introspection_protocol"
	"github.com/google/uuid"
)

type CompilerIntrospectionServer struct {
	*IntrospectionServer
	Quit chan bool
}

func NewCompilerIntrospectionServer() (*CompilerIntrospectionServer, error) {
	introspection_protocol.RegisterTypes()

	server, err := NewServer()
	if err != nil {
		return nil, err
	}

	return &CompilerIntrospectionServer{
		IntrospectionServer: server,
		Quit:                make(chan bool),
	}, nil
}

func (s *CompilerIntrospectionServer) WaitForClient() {
	for {
		select {
		case req := <-s.RequestChan:
			switch request := req.(type) {
			case introspection_protocol.HeloRequest:
				clientID, err := s.registerClient(request.IntrospectionType)

				if err != nil {
					s.ResponseChan <- err
				} else {
					response := introspection_protocol.HeloResponse{
						ClientID: clientID,
					}
					s.ResponseChan <- response
					return
				}
			default:
				s.ResponseChan <- errors.New("Expected HELO request")
			}
		case <-s.Quit:
			return
		}
	}
}

func (s *CompilerIntrospectionServer) registerClient(introspectionType introspection_protocol.IntrospectionType) (string, error) {
	newClientID := uuid.New().String()
	return newClientID, nil
}

func (s *CompilerIntrospectionServer) Shutdown() {
	s.Quit <- true
}
