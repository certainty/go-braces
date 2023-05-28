package frame

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	termHeight := m.height
	termWidth := m.width

	topBarView := m.topBar.View()
	topBarHeight := lipgloss.Height(topBarView)

	eventLogView := m.eventLog.View()
	eventlogHeight := lipgloss.Height(eventLogView)

	statusBarView := m.statusBar.View()
	statusBarHeight := lipgloss.Height(statusBarView)

	mainView := m.mainContainerStyle.
		Copy().
		Width(termWidth).
		Height(termHeight - topBarHeight - eventlogHeight - statusBarHeight).
		Render("Main container")

	return m.frameStyle.
		Copy().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, topBarView, mainView, eventLogView, statusBarView))
}
