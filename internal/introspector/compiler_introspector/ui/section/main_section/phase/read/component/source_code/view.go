package source_code

func (m Model) View() string {
	content := m.codeContentStyle.
		Copy().
		Width(m.codeViewer.Width - 2).
		Height(m.codeViewer.Height - 2).
		Render(m.code)

	m.codeViewer.SetContent(content)
	return m.codeViewer.View()
}
