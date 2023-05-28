package frame

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	termHeight := m.height
	termWidth := m.width

	topBarView := m.topBar.View()
	mainSectionView := m.mainSection.View()
	statusBarView := m.statusBar.View()

	views := []string{topBarView, mainSectionView}
	if m.eventLog.IsVisible() {
		views = append(views, m.eventLog.View())
	}
	views = append(views, statusBarView)

	return m.frameStyle.
		Copy().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, views...))
}
