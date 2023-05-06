package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/introspection/service/compiler_introsection"
	"github.com/certainty/go-braces/internal/repl"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/spf13/cobra"
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

	if enableCompilerIntrospection || enableVMIntrospection {
		var wg sync.WaitGroup

		if enableCompilerIntrospection {
			wg.Add(1)

			compilerIntrospectionServer, err := compiler_introsection.StartServer(&wg)

			if err != nil {
				log.Fatal("Could not start introspection server: ", err.Error())
				os.Exit(1)
			}
			runRepl(compilerIntrospectionServer)
			compilerIntrospectionServer.Stop()
			wg.Wait()
		}
	} else {
		runRepl(nil)
	}
}

func runRepl(compilerIntrospection *introspection.IntrospectionServer) {
	compilerOptions := compiler.DefaultOptions()

	if compilerIntrospection != nil {
		compilerOptions = compiler.NewCompilerOptions(compilerIntrospection.Api)
	}

	vm := vm.NewVM(vm.DefaultOptions())
	compiler := compiler.NewCompiler(compilerOptions)

	repl, err := repl.NewRepl(vm, compiler, compilerIntrospection)
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
