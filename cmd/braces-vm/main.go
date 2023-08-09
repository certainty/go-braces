package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)

	compiler_introspection.RegisterTypes()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
