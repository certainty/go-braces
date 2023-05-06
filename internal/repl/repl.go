package repl

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/chzyer/readline"
)

type Repl struct {
	vm                    *vm.VM
	compiler              *compiler.Compiler
	lineedit              *readline.Instance
	compilerIntrospection *introspection.IntrospectionServer
	inputCount            int
}

func NewRepl(vm *vm.VM, compiler *compiler.Compiler, compilerIntrospection *introspection.IntrospectionServer) (*Repl, error) {
	rl, err := readline.New("> ")
	if err != nil {
		return nil, err
	}

	return &Repl{
		vm:                    vm,
		compiler:              compiler,
		compilerIntrospection: compilerIntrospection,
		lineedit:              rl,
		inputCount:            0,
	}, nil
}

// run without introspection
func (r *Repl) Run() {
	println("Welcome to the Go Braces REPL!")
	println("Type :exit or CTRL-C for exit, and :help for help")

	if r.compilerIntrospection != nil {
		println("Compiler Introspection is enabled. To connect run: braces-introspect compiler", r.compilerIntrospection.ListenAddr)
		println("Waiting for introspection client ....")
		r.compilerIntrospection.WaitForClient()
	}

	println("\n")

	defer r.lineedit.Close()

	for {
		input, err := r.getInput()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		wasCommand, doExit, err := r.handleCommand(input)
		if err != nil {
			fmt.Println(err.Error())
		}

		if wasCommand {
			if doExit {
				break
			}
		} else {
			result, err := r.compileAndRun(input)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("=> %v\n", result)
			}
		}
	}

}

func (r *Repl) handleCommand(input string) (bool, bool, error) {
	if input == ":help" {
		fmt.Println("Show help")
		return true, false, nil
	}

	if input == ":exit" {
		return true, true, nil
	}
	return false, false, nil
}

func (r *Repl) getInput() (string, error) {
	openBrackes, openParens, openQuotes := 0, 0, 0
	buffer := ""
	r.lineedit.SetPrompt(fmt.Sprintf("%03d:> ", r.inputCount))

	for {
		line, err := r.lineedit.Readline()
		if err != nil {
			return "", err
		}

		buffer += line

		for _, char := range line {
			switch char {
			case '[':
				openBrackes++
			case ']':
				openBrackes--
			case '(':
				openParens++
			case ')':
				openParens--
			case '"':
				openQuotes++
			}
		}

		if openBrackes == 0 && openParens == 0 && openQuotes%2 == 0 {
			r.inputCount += 1
			r.lineedit.SetPrompt(fmt.Sprintf("%03d:> ", r.inputCount))
			return buffer, nil
		} else {
			buffer += "\n"
			r.lineedit.SetPrompt(fmt.Sprintf("%03d:* ", r.inputCount))
		}
	}
}

func (r *Repl) compileAndRun(input string) (isa.Value, error) {
	assemblyModule, err := r.compiler.CompileString(input)

	if err != nil {
		return nil, err
	}
	return r.vm.ExecuteModule(assemblyModule)
}
