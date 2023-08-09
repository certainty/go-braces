package main

import (
	"fmt"
	"github.com/certainty/go-braces/pkg/vm/run"
	"github.com/spf13/cobra"
	"os"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Compile and run the file provided",
	Long:  `More documentation still to come here`,
	Run: func(cmd *cobra.Command, args []string) {
		runInput(cmd, args)
	},
}

func runInput(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		println("Please provide a file to run")
		return
	}

	enableCompilerIntrospection, _ := cmd.Flags().GetBool("introspect-compiler")
	enableVMIntrospection, _ := cmd.Flags().GetBool("introspect-vm")

	options := run.Options{
		IntrospectCompiler: enableCompilerIntrospection,
		IntrospectVM:       enableVMIntrospection,
	}

	runner, err := run.NewRunner(options)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	runner.Run(args[0])
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("introspect-compiler", "c", false, "Run the repl in introspection mode for the compiler")
	runCmd.Flags().BoolP("introspect-vm", "m", false, "Run the repl in introspection mode for the vm")
}
