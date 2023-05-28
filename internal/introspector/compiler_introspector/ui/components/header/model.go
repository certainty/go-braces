package header

// import (
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type Model struct {
// 	containerWidth  int
// 	containerHeight int
// 	theme           theme.Theme
// 	Title           string
// 	Connected       bool
// 	RequestState    components.IntrospectionRequestState
// }

// func NewModel(theme theme.Theme, title string, connected bool) Model {
// 	return Model{
// 		theme:        theme,
// 		Title:        title,
// 		Connected:    connected,
// 		RequestState: components.NoRequest,
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	return m, nil
// }

// func (m Model) View() string {
// 	titleWidth := 50
// 	restWidth := m.containerWidth - titleWidth
// 	title := m.theme.Title.Width(titleWidth).Render(m.Title)
// 	status := "disconnected"
// 	if m.Connected {
// 		status = "connected"
// 	}

// 	connectionStatus := ""
// 	if m.RequestState == components.NoRequest {
// 		connectionStatus = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("\uF1E6 " + status)
// 	} else {
// 		connectionStatus = m.theme.Header.Width(restWidth).Align(lipgloss.Right).PaddingRight(2).Render("⇄ " + status)
// 	}

// 	return lipgloss.JoinHorizontal(lipgloss.Top, title, connectionStatus)
// }

// func (m Model) Resize(width, height int) components.Model {
// 	m.containerWidth = width
// 	m.containerHeight = height
// 	return m
// }

// func (m Model) ContainerHeight() int {
// 	return m.containerHeight
// }

// func (m Model) ContainerWidth() int {
// 	return m.containerWidth
// }
