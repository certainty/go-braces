package phase_indicator

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerWidth = msg.Width
		m.containerHeight = msg.Height
	case MsgReset:
		m.currentPhase = 0
		m.finished = false
	case MsgFinish:
		m.finished = true
	case MsgPhase:
		for i, p := range m.phases {
			if compiler_introspection.CompilationPhase(msg) == p {
				m.currentPhase = i
				break
			}
		}
	}
	return m, nil
}
