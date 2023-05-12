package main

import (
	"fmt"
	"log"
	"os"

	"github.com/certainty/go-braces/internal/introspection"
)

func main() {
	introspection.RegisterWireTypes()
	wireClient := introspection.NewWireClient("compiler")
	wireControlConnection, _, err := wireClient.Connect()
	if err != nil {
		log.Fatalf("could not connect to server %v", err)
	}

	wireControlConnection.Out <- introspection.Hello{}
	<-wireControlConnection.In
	// read from stdin
	intro := make([]byte, 1)
	if _, err := os.Stdin.Read(intro); err != nil {
		fmt.Println(err)
	}
}
