package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "braces-compile",
	Short: "A compiler for the braces language, with emphasis on introspection.",
	Long:  `A compiler for the braces language, with emphasis on introspection.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Add your application logic here
		fmt.Println("Runnning braces-compile")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
