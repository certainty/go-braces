package main_section

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
	}

	return m, nil
}
