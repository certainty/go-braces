package repl

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/internal/vm"
	"github.com/certainty/go-braces/internal/vm/language/value"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/bubbline/editline"
)

type Repl struct {
	editline *editline.Model
	vm       *vm.VM
	lines    []string
}

func NewRepl(width, height int, vm *vm.VM) *Repl {
	return &Repl{
		editline: editline.New(width, height),
		vm:       vm,
		lines:    make([]string, 0),
	}
}

func (r *Repl) Init() tea.Cmd {
	return r.editline.Init()
}

func (r *Repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// handle evaluations first
	switch msg := msg.(type) {
	case ExecutionFailed:
		{
			// TODO: format error
			r.lines = append(r.lines, msg.Error.Error())
		}
	case ExecutionSucceeded:
		{
			r.lines = append(r.lines, fmt.Sprintf("%v", msg.Value))
		}
	}

	// now handle input
	m, cmd := r.editline.Update(msg)
	r.editline = m.(*editline.Model)
	switch msg.(type) {
	case editline.InputCompleteMsg:
		cmd = r.RunVMCmd()
	}
	return r, cmd
}

func (r *Repl) View() string {
	// return all lines joined by newline and finally the editline view
	return strings.Join(r.lines, "\n") + "\n" + r.editline.View() + "\n"
}

type ExecutionFailed struct {
	Error error
}
type ExecutionSucceeded struct {
	Value value.Value
}

func (r *Repl) RunVMCmd() tea.Cmd {
	return func() tea.Msg {
		return ExecutionFailed{
			Error: fmt.Errorf("not implemented"),
		}
	}
}
