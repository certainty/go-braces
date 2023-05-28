package compilation_info

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m Model) View() string {
	options := "N/A"
	if len(m.CompilerOptions) > 0 {
		options = strings.Join(m.CompilerOptions, " ")
	}

	inputDescription := "N/A"
	if m.Origin != nil {
		inputDescription = (*m.Origin).Description()
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.theme.Info.Copy().PaddingLeft(1).PaddingRight(3).Render("Input: "+inputDescription),
		m.theme.Info.Copy().Render("Options: "+options),
	)

	width := m.containerWidth
	height := m.containerHeight - 2

	return lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder(), true, false, true, false).
		BorderForeground(m.theme.Colors.InactiveBorder).
		Width(width).
		Height(height).
		Render(content)
}
