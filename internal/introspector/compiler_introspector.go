package introspector

import (
	"context"
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/introspection/service/compiler_introsection"
	"google.golang.org/grpc"
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

	// run a first helo
	capability := compiler_introsection.Capability{}
	heloResponse, err := client.Helo(context.Background(), &capability)
	if err != nil {
		log.Fatalf("could not call Helo: %v", err)
	}
	fmt.Printf("Helo response: %v\n", heloResponse)

	return nil
}
