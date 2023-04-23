package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "braces-vm",
	Short: "A virtual machine for the Braces language, with emphasis on introspection.",
	Long: `This is the virual machine for the braces language. 
          It runs pre compiled (ahead of time compiled) braces bytecode. 
          The VM has been built to help understanding what it does during runtime and thus provides
          a set of tools and options to do so.
        `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Runnning braces-vm")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
