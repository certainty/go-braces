package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspection"
)

func main() {
	introspection.RegisterTypes()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
