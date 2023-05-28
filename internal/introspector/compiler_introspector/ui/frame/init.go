package frame

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	cmds = append(cmds, CmdConnectClient(m.client))
	cmds = append(cmds, common.CmdTick())
	cmds = append(cmds, m.statusBar.Init())

	return tea.Batch(cmds...)
}
