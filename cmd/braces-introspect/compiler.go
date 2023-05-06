package main

import (
	"github.com/certainty/go-braces/internal/introspector"
	"github.com/spf13/cobra"
)

var compilerCmd = &cobra.Command{
	Use:   "compiler SERVER_ADDRESS",
	Short: "Start an introspector for the braces compiler",
	Long:  `More documentation still to come here`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func run(cmd *cobra.Command, args []string) {
	serverAddress := args[0]
	introspectionClient := introspector.NewCompilerIntrospector(serverAddress)
	introspectionClient.Start()
}

func init() {
	rootCmd.AddCommand(compilerCmd)
}
