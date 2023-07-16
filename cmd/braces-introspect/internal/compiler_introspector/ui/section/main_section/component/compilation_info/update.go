package compilation_info

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
	case MsgNewCompilation:
		m.Origin = &msg.Origin
	default:
		break
	}

	return m, nil
}
