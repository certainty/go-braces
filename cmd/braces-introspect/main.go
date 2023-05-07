package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspection/introspection_events"
)

func main() {
	introspection.RegisterTypes()
	introspection_events.RegisterEventTypes()
	compiler_introspection.RegisterControlTypes()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
