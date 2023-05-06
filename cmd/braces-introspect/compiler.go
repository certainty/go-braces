package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/introspector"
	"github.com/spf13/cobra"
)

var compilerCmd = &cobra.Command{
	Use:   "compiler IPC_DIR",
	Short: "Start an introspector for the braces compiler",
	Long:  `More documentation still to come here`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func run(cmd *cobra.Command, args []string) {
	ipcDir := args[0]
	if err := introspector.RunIntrospector(ipcDir); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(compilerCmd)
}
