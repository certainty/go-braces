package repl

import (
	"errors"
	"fmt"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/knz/bubbline"
	"io"
)

type Repl struct {
	lineEditor *bubbline.Editor
	vm         *vm.VM
}

func NewRepl(vm *vm.VM) *Repl {
	return &Repl{
		lineEditor: bubbline.New(),
		vm:         vm,
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
		println(val)
	}
}
