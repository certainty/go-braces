package repl

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/chzyer/readline"
)

type Repl struct {
	vmInstance            *vm.VM
	compiler              *compiler.Compiler
	lineedit              *readline.Instance
	compilerIntrospection *compiler_introspection.Server
	inputCount            int
}

type Options struct {
	IntrospectCompiler bool
	IntrospectVM       bool
}

func NewRepl(options Options) (*Repl, error) {
	rl, err := readline.New("> ")
	if err != nil {
		return nil, err
	}

	var compiler *compiler.Compiler
	var compilerIntrospection *compiler_introspection.Server

	if options.IntrospectCompiler {
		compiler, compilerIntrospection, err = newCompilerWithIntrospection()
		if err != nil {
			return nil, err
		}
	} else {
		compiler = newCompiler()
	}

	instance := vm.NewVM(vm.DefaultOptions())
	if err != nil {
		return nil, err
	}

	return &Repl{
		vmInstance:            instance,
		compiler:              compiler,
		compilerIntrospection: compilerIntrospection,
		lineedit:              rl,
		inputCount:            0,
	}, nil
}

func newCompiler() *compiler.Compiler {
	return compiler.NewCompiler(compiler.DefaultOptions())
}

func newCompilerWithIntrospection() (*compiler.Compiler, *compiler_introspection.Server, error) {
	introspectionServer, err := compiler_introspection.NewServer()

	if err != nil {
		return nil, nil, err
	}
	compilerOptions := compiler.NewCompilerOptions(introspectionServer.API())
	return compiler.NewCompiler(compilerOptions), introspectionServer, nil
}

// run without introspection
func (r *Repl) Run() {
	println("Welcome to the Go Braces REPL!")
	println("Type :exit or CTRL-C for exit, and :help for help")

	if r.compilerIntrospection != nil {
		println("Compiler Introspection is enabled. To connect run: braces-introspect compiler", r.compilerIntrospection.IPCDir())
		println("Waiting for introspection client ....")
		clientId := r.compilerIntrospection.WaitForClient()
		println("Client connected. ID is", clientId)
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
				fmt.Printf("%s\n", r.vmInstance.WriteValue(result))
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
	return r.vmInstance.ExecuteModule(assemblyModule)
}
