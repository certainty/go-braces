package compilation_info

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
	case common.MsgIntrospectionEvent:
		switch msg := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			m.Origin = &msg.Origin
		default:
			break
		}
	default:
		break
	}

	return m, nil
}
