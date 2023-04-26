package repl

import (
	"errors"
	"fmt"
	"io"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/certainty/go-braces/internal/vm/language/value"
	"github.com/knz/bubbline"
)

type Repl struct {
	lineEditor *bubbline.Editor
	vm         *vm.VM
	compiler   *compiler.Compiler
}

func NewRepl(vm *vm.VM, compiler *compiler.Compiler) *Repl {
	return &Repl{
		lineEditor: bubbline.New(),
		vm:         vm,
		compiler:   compiler,
	}
}

// run without introspection
func (r *Repl) Run() {
	println("Welcome to the Go Braces REPL!")

	for {
		val, err := r.lineEditor.GetLine()

		if err != nil {
			if err == io.EOF {
				// No more input.
				break
			}
			if errors.Is(err, bubbline.ErrInterrupted) {
				// Entered Ctrl+C to cancel input.
				fmt.Println("^C")
			} else if errors.Is(err, bubbline.ErrTerminated) {
				fmt.Println("terminated")
				break
			} else {
				fmt.Println("error:", err)
			}
			continue
		}
		r.lineEditor.AddHistory(val)

		// TODO: check if input is a command or normal input
		result, err := r.compileAndRun(val)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func (r *Repl) compileAndRun(input string) (value.Value, error) {
	complicationUnit, err := r.compiler.JitCompile(input)

	if err != nil {
		return nil, err
	}
	return r.vm.Execute(complicationUnit)
}
