package source_code

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	content := m.codeContentStyle.
		Copy().
		Width(m.codeViewer.Width - 2).
		Height(m.codeViewer.Height - 2).
		Render(m.code)

	m.codeViewer.SetContent(content)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		" Source",
		m.codeViewer.View(),
	)
}
