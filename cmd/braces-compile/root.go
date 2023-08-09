package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "braces-compile",
	Short: "A compiler for the braces language, with emphasis on introspection.",
	Long:  `A compiler for the braces language, with emphasis on introspection.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("braces-compile called %v", args)
		if len(args) == 0 {
			fmt.Println("No input file specified")
			os.Exit(1)
		}
	},
}

func init() {
}
