package topbar

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerWidth = msg.Width
		m.containerHeight = msg.Height
	case common.MsgClientConnected:
		m.isConnected = bool(msg)
	}
	return m, nil
}
