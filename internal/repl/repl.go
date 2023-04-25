package repl

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/knz/bubbline/editline"
)

type Repl struct {
	editline *editline.Model
}

func NewRepl(width, height int) *Repl {
	return &Repl{
		editline: editline.New(width, height),
	}
}

func (r *Repl) Init() tea.Cmd {
	return r.editline.Init()
}

func (r *Repl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m, cmd := r.editline.Update(msg)
	r.editline = m.(*editline.Model)
	return r, cmd
}

func (r *Repl) View() string {
	return r.editline.View()
}
