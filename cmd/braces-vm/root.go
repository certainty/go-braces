package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-braces.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
