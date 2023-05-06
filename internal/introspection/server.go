package introspection

import "context"

type IntrospectionServer struct {
	Api        API
	ListenAddr string
	Context    context.Context
	Cancel     context.CancelFunc
}

func NewIntrospectionServer(addr string, ctx context.Context, cancel context.CancelFunc) *IntrospectionServer {
	return &IntrospectionServer{
		Api:        NewAPI(),
		ListenAddr: addr,
		Context:    ctx,
		Cancel:     cancel,
	}
}

func (s *IntrospectionServer) Stop() {
	s.Cancel()
}
