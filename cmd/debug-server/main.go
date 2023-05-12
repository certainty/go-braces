package main

import (
	"log"

	"github.com/certainty/go-braces/internal/introspection"
)

func main() {
	introspection.RegisterWireTypes()
	wireServer, err := introspection.NewWireServer("compiler")
	if err != nil {
		log.Fatalf("could not create server %v", err)
	}

	wireControlConnection, wireEventConnection, err := wireServer.Accept()
	if err != nil {
		log.Fatalf("could not accept connection %v", err)
	}

	req := <-wireControlConnection.In
	log.Printf("got request %v", req)

	wireControlConnection.Close()
	wireEventConnection.Close()
}
