package topbar

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	titleWidth := 50
	restWidth := m.containerWidth - titleWidth
	title := m.theme.Title.Width(titleWidth).Render(m.Title)
	status := "disconnected"
	if m.isConnected {
		status = "connected"
	}

	requestIndicator := ""
	if m.requestStatus == common.NoRequest {
		requestIndicator = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("\uF1E6 " + status)
	} else {
		requestIndicator = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("⇄ " + status)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, requestIndicator)
}
