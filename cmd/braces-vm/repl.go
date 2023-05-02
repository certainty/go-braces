package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/repl"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/spf13/cobra"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Start the VM in repl mode",
	Long:  `More documentation still to come here`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	vm := vm.NewVM(vm.DefaultOptions())
	compiler := compiler.NewCompiler(compiler.DefaultOptions())

	repl, err := repl.NewRepl(vm, compiler)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	repl.Run()
}

func init() {
	rootCmd.AddCommand(replCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
