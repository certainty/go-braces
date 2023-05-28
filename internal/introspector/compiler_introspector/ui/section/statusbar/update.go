package statusbar

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.RequestSpinner, cmd = m.RequestSpinner.Update(msg)
		return m, cmd
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
	case common.MsgModeChange:
		m.Mode = msg.ActiveMode
	case common.MsgClientConnected:
		m.isConnected = bool(msg)
	case common.MsgError:
		m.err = msg.Err
	}
	return m, nil
}
