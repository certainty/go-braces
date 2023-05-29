package frame

import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
	termHeight := m.height
	termWidth := m.width

	topBarView := m.sectionTopBar.View()
	mainSectionView := m.sectionMain.View()
	statusBarView := m.sectionStatusBar.View()

	views := []string{topBarView, mainSectionView}
	if m.sectionEventLog.IsVisible() {
		views = append(views, m.sectionEventLog.View())
	}
	views = append(views, statusBarView)

	return m.styleFrame.
		Copy().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, views...))
}
