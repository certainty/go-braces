package frame

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/eventlog"
	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	termHeight := m.height
	termWidth := m.width

	topBarView := m.sections[SectionTopBar].View()
	mainSectionView := m.sections[SectionMain].View()
	statusBarView := m.sections[SectionStatusBar].View()

	views := []string{topBarView, mainSectionView}
	if m.sections[SectionEventLog].(eventlog.Model).IsVisible() {
		views = append(views, m.sections[SectionEventLog].View())
	}
	views = append(views, statusBarView)

	return m.styleFrame.
		Copy().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, views...))
}
