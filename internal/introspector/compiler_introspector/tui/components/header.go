package components

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HeaderModel struct {
	ContainerWidth  int
	ContainerHeight int
	theme           theme.Theme
	Title           string
	Connected       bool
	RequestState    IntrospectionRequestState
}

func InitialHeaderModel(theme theme.Theme, title string, connected bool) HeaderModel {
	return HeaderModel{
		theme:        theme,
		Title:        title,
		Connected:    connected,
		RequestState: NoRequest,
	}
}

func (m HeaderModel) Init() tea.Cmd {
	return nil
}

func (m HeaderModel) Update(msg tea.Msg) (HeaderModel, tea.Cmd) {
	return m, nil
}

func (m HeaderModel) View() string {
	titleWidth := 50
	restWidth := m.ContainerWidth - titleWidth
	title := m.theme.Title.Width(titleWidth).Render(m.Title)
	status := "disconnected"
	if m.Connected {
		status = "connected"
	}

	connectionStatus := ""
	if m.RequestState == NoRequest {
		connectionStatus = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("\uF1E6 " + status)
	} else {
		connectionStatus = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("⇄ " + status)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, title, connectionStatus)
}
