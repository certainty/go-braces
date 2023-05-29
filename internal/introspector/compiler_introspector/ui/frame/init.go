package frame

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(common.CmdTick(), CmdConnectClient(m.client), m.sections[SectionStatusBar].Init())
}
