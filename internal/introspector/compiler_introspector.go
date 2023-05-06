package introspector

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/certainty/go-braces/internal/introspection/service/compiler_introsection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type CompilerIntrospector struct {
	address string
}

func NewCompilerIntrospector(address string) *CompilerIntrospector {
	return &CompilerIntrospector{address: address}
}

func (introspector *CompilerIntrospector) Start() error {
	conn, err := grpc.Dial(introspector.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := compiler_introsection.NewCompilerIntrospectionClient(conn)

	capabilities := compiler_introsection.CapabilityList{}
	request := compiler_introsection.HeloRequest{
		Capabilities:      &capabilities,
		IntrospectionType: compiler_introsection.IntrospectionType_COMPILER,
	}

	heloResponse, err := client.Helo(context.Background(), &request)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("gRPC error: %s (code: %s)\n", s.Message(), s.Code())
		} else {
			log.Fatalf("could not call Helo: %v", err)
		}
	}
	fmt.Printf("Helo response: %v\n", heloResponse)

	eventStream, err := client.EventStream(context.Background(), &compiler_introsection.EventStreamRequest{})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			fmt.Printf("gRPC error: %s (code: %s)\n", s.Message(), s.Code())
		} else {
			log.Fatalf("could not call EventStream: %v", err)
		}
	}

	// Process events from the server
	for {
		event, err := eventStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive event: %v", err)
		}
		fmt.Printf("Received event: %s\n", event.Json)
	}

	return nil
}
