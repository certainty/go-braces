package statusbar

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Init() tea.Cmd {
	return m.RequestSpinner.Tick
}
