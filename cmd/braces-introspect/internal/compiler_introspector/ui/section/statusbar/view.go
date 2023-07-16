package statusbar

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) RenderShortCut(shortcut key.Binding) string {
	return m.theme.Statusbar.Copy().Faint(!shortcut.Enabled()).PaddingRight(2).Render(shortcut.Help().Key)
}

func (m Model) View() string {
	renderedShortcuts := []string{
		"\uF11C  ",
	}

	if len(m.contextKeyMap.ShortCuts) > 0 {
		for _, shortcut := range m.contextKeyMap.ShortCuts {
			renderedShortcuts = append(renderedShortcuts, m.RenderShortCut(shortcut))
		}

		renderedShortcuts = append(renderedShortcuts, m.theme.Statusbar.Render("|"))
	}

	for _, shortcut := range m.globalKeyMap.ShortCuts {
		renderedShortcuts = append(renderedShortcuts, m.RenderShortCut(shortcut))
	}

	var statusMessage string
	if m.err != nil {
		statusMessage = m.theme.Statusbar.Copy().Foreground(m.theme.Colors.Error).Render(m.err.Error())
	} else {
		if m.IsConnected {
			statusMessage = m.theme.Statusbar.Copy().Foreground(m.theme.Colors.Success).Render("Connected to introspection server")
		} else {
			statusMessage = m.theme.Statusbar.Copy().Render("No connection to introspection server")
		}
	}

	connectionStatus := ""
	if m.RequestState != common.NoRequest {
		connectionStatus = m.RequestSpinner.View()
	}

	m.impl.SetSize(m.containerWidth)
	m.impl.SetContent(m.Mode.String(), statusMessage, lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...), connectionStatus)

	return m.impl.View()
}
