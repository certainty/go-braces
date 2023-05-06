package introspection

import (
	"context"

	"github.com/certainty/go-braces/internal/introspection/service"
)

type IntrospectionServer struct {
	Api           API
	ListenAddr    string
	ServiceEvents chan service.Event
	Context       context.Context
	Cancel        context.CancelFunc
}

func NewIntrospectionServer(addr string, ctx context.Context, cancel context.CancelFunc) *IntrospectionServer {
	return &IntrospectionServer{
		Api:           NewAPI(),
		ListenAddr:    addr,
		Context:       ctx,
		Cancel:        cancel,
		ServiceEvents: make(chan service.Event),
	}
}

func (s *IntrospectionServer) Stop() {
	s.Cancel()
}

func (s *IntrospectionServer) WaitForClient() error {
	for {
		evt := <-s.ServiceEvents
		switch evt.(type) {
		case service.ClientConnected:
			return nil
		default:
			continue
		}
	}
}
