package compiler_introsection

import (
	context "context"
	"encoding/json"
	"log"
	"net"
	"sync"

	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/introspection/service"
	"github.com/google/uuid"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type server struct {
	UnimplementedCompilerIntrospectionServer
	api     introspection.API
	service chan service.Event
}

func StartServer(wg *sync.WaitGroup) (*introspection.IntrospectionServer, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	introspectionServer := introspection.NewIntrospectionServer(lis.Addr().String(), ctx, cancel)

	grpcServer := grpc.NewServer()
	RegisterCompilerIntrospectionServer(grpcServer, &server{
		api:     introspectionServer.Api,
		service: introspectionServer.ServiceEvents,
	})

	go func() {
		defer wg.Done()

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve commpiler introspection: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return introspectionServer, nil
}

func (s *server) Helo(ctx context.Context, req *HeloRequest) (*HeloResponse, error) {
	newClientID := uuid.New().String()
	s.service <- service.ClientConnected{ClientID: newClientID}

	if req.IntrospectionType != IntrospectionType_COMPILER {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid introspection type. Expected COMPILER but got: %v", req.IntrospectionType)
	}

	return &HeloResponse{
		ClientId:     newClientID,
		Capabilities: &CapabilityList{},
	}, nil
}

func (s *server) EventStream(req *EventStreamRequest, stream CompilerIntrospection_EventStreamServer) error {
	for {
		// TODO: in  the future I might want to select here
		evt := s.api.ReceiveEvent()
		eventData, err := serializeEvent(evt)

		if err != nil {
			continue
		}
		rpcEvent := &Event{
			Json: eventData,
		}

		if err := stream.Send(rpcEvent); err != nil {
			return err
		}
	}
}

func serializeEvent(evt introspection.IntrospectionEvent) (string, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *server) AbortSession(ctx context.Context, req *AbortSessionRequest) (*AbortSessionResponse, error) {
	return nil, nil
}
