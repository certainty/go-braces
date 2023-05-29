package read

func (m Model) View() string {
	sourceCodeView := m.sourceCodePane.View()
	return sourceCodeView
}
