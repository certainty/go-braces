package repl

import (
	"fmt"
	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/vm"
	"github.com/certainty/go-braces/internal/vm/language/value"
	"github.com/chzyer/readline"
)

type Repl struct {
	vm         *vm.VM
	compiler   *compiler.Compiler
	lineedit   *readline.Instance
	inputCount int
}

func NewRepl(vm *vm.VM, compiler *compiler.Compiler) (*Repl, error) {
	rl, err := readline.New("> ")
	if err != nil {
		return nil, err
	}

	return &Repl{
		vm:         vm,
		compiler:   compiler,
		lineedit:   rl,
		inputCount: 0,
	}, nil
}

// run without introspection
func (r *Repl) Run() {
	println("Welcome to the Go Braces REPL!")
	println("Press Ctrl+C to exit and :help for help\n")
	defer r.lineedit.Close()

	for {
		input, err := r.getInput()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		result, err := r.compileAndRun(input)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("=> %v\n", result)
		}
	}
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

func (r *Repl) compileAndRun(input string) (value.Value, error) {
	complicationUnit, err := r.compiler.JitCompile(input)

	if err != nil {
		return nil, err
	}
	return r.vm.Execute(complicationUnit)
}
