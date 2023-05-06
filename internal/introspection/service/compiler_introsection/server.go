package compiler_introsection

import (
	context "context"
	"log"
	"net"
	"sync"

	"github.com/certainty/go-braces/internal/introspection"
	grpc "google.golang.org/grpc"
)

type server struct {
	UnimplementedCompilerIntrospectionServer
	api introspection.API
}

// starte the server on the next available port in a separate goroutine
func StartServer(ctx context.Context, wg *sync.WaitGroup, api introspection.API) (net.Addr, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()
	RegisterCompilerIntrospectionServer(grpcServer, &server{
		api: api,
	})

	// Start the server in a separate goroutine
	go func() {
		defer wg.Done()

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return lis.Addr(), nil
}

func (s *server) Hello(ctx context.Context, capability Capability) (Capability, error) {
	return Capability{}, nil
}

func (s *server) StartSession(ctx context.Context, req *StartSessionRequest) (*StartSessionResponse, error) {
	return nil, nil
}

func (s *server) EventStream(req *StartSessionResponse, stream CompilerIntrospection_EventStreamServer) error {
	return nil
}

func (s *server) AbortSession(ctx context.Context, req *AbortSessionRequest) (*AbortSessionResponse, error) {
	return nil, nil
}
