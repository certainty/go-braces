package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector"
	"github.com/spf13/cobra"
)

var compilerCmd = &cobra.Command{
	Use:   "compiler [flags]",
	Short: "Start an introspector for the braces compiler",
	Long:  `More documentation still to come here`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func run(cmd *cobra.Command, args []string) {
	if err := compiler_introspector.RunIntrospector(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(compilerCmd)
}
