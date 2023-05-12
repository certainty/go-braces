package main

import (
	"fmt"
	"github.com/certainty/go-braces/internal/repl"
	"github.com/spf13/cobra"
	"os"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start the VM in repl mode",
	Long:  `More documentation still to come here`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

func run(cmd *cobra.Command, args []string) {
	enableCompilerIntrospection, _ := cmd.Flags().GetBool("introspect-compiler")
	enableVMIntrospection, _ := cmd.Flags().GetBool("introspect-vm")

	options := repl.Options{
		IntrospectCompiler: enableCompilerIntrospection,
		IntrospectVM:       enableVMIntrospection,
	}

	repl, err := repl.NewRepl(options)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	repl.Run()
}

func init() {
	rootCmd.AddCommand(replCmd)

	replCmd.Flags().BoolP("introspect-compiler", "c", false, "Run the repl in introspection mode for the compiler")
	replCmd.Flags().BoolP("introspect-vm", "m", false, "Run the repl in introspection mode for the vm")
}
