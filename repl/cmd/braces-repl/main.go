package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "braces-debug",
	Short: "A repl for the braces language.",
	Long:  `A repl for the braces language.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Runnning braces-repl")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
