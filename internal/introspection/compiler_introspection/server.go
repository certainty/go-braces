package compiler_introspection

import (
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/google/uuid"
)

type Server struct {
	*introspection.Server
	Quit chan bool
}

func NewServer() (*Server, error) {
	server, err := introspection.NewServer()
	if err != nil {
		return nil, err
	}

	return &Server{
		server,
		make(chan bool),
	}, nil
}

func (s *Server) WaitForClient() string {
	for {
		select {
		case req := <-s.RequestChan:
			switch request := req.(type) {
			case introspection.HeloRequest:
				clientID, err := s.registerClient(request.IntrospectionType)

				if err != nil {
					s.ResponseChan <- err
				} else {
					response := introspection.HeloResponse{
						ClientID: clientID,
					}
					s.ResponseChan <- response
					return clientID
				}
			default:
				s.ResponseChan <- "Expected HELO request"
			}
		case <-s.Quit:
			return ""
		}
	}
}

func (s *Server) registerClient(introspectionType introspection.IntrospectionType) (string, error) {
	newClientID := uuid.New().String()
	return newClientID, nil
}

func (s *Server) Shutdown() {
	s.Quit <- true
}

func (s *Server) API() introspection.API {
	return nil
}
