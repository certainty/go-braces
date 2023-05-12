package repl

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/compiler/input"
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
	options               Options
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

	instance := vm.NewVM(vm.DefaultOptions())
	if err != nil {
		return nil, err
	}

	return &Repl{
		vmInstance: instance,
		lineedit:   rl,
		inputCount: 0,
		options:    options,
	}, nil
}

// run without introspection
func (r *Repl) Run() {
	println("Welcome to the Go Braces REPL!")
	println("Type :exit or CTRL-C for exit, and :help for help")

	if r.options.IntrospectCompiler {
		println("Compiler Introspection is enabled")
		println("Waiting for introspection client ....")

		var err error
		r.compilerIntrospection, err = compiler_introspection.NewServer()
		if err != nil {
			return
		}

		instrumenter, err := r.compilerIntrospection.WaitForClient()
		if err != nil {
			return
		}

		r.compiler = compiler.NewCompiler(compiler.NewCompilerOptions(instrumenter))
		r.doRun()
	} else {

		r.compiler = compiler.NewCompiler(compiler.DefaultOptions())
		r.doRun()
	}
}

func (r *Repl) doRun() {
	defer r.lineedit.Close()
	for {

		// now get input
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

func (r *Repl) compileAndRun(source string) (isa.Value, error) {
	replInput := input.NewReplInput(uint64(r.inputCount), source)
	assemblyModule, err := r.compiler.CompileModule(replInput)

	if err != nil {
		return nil, err
	}
	return r.vmInstance.ExecuteModule(assemblyModule)
}
