package compilation_info

// import (
// 	"strings"

// 	"github.com/certainty/go-braces/internal/compiler/input"
// 	"github.com/certainty/go-braces/internal/compiler/location"
// 	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type Model struct {
// 	containerWidth  int
// 	containerHeight int
// 	theme           theme.Theme
// 	Origin          *location.Origin
// 	CompilerOptions []string
// }

// func NewModel(theme theme.Theme, input *input.Input, options []string) Model {
// 	return Model{
// 		theme:           theme,
// 		CompilerOptions: options,
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case messages.IntrospectionEventMsg:
// 		switch msg := msg.Event.(type) {
// 		case compiler_introspection.EventBeginCompileModule:
// 			m.Origin = &msg.Origin
// 		default:
// 			break
// 		}
// 	default:
// 		break
// 	}

// 	return m, nil
// }

// func (m Model) View() string {
// 	options := "N/A"
// 	if len(m.CompilerOptions) > 0 {
// 		options = strings.Join(m.CompilerOptions, " ")
// 	}

// 	inputDescription := "N/A"
// 	if m.Origin != nil {
// 		inputDescription = (*m.Origin).Description()
// 	}

// 	content := lipgloss.JoinHorizontal(
// 		lipgloss.Top,
// 		m.theme.Info.Copy().PaddingLeft(1).PaddingRight(3).Render("Input: "+inputDescription),
// 		m.theme.Info.Copy().Render("Options: "+options),
// 	)

// 	width := m.containerWidth
// 	height := m.containerHeight - 2

// 	return lipgloss.
// 		NewStyle().
// 		Border(lipgloss.RoundedBorder(), true, false, true, false).
// 		BorderForeground(m.theme.Colors.InactiveBorder).
// 		Width(width).
// 		Height(height).
// 		Render(content)
// }

// func (m Model) Resize(width, height int) components.Model {
// 	m.containerWidth = width
// 	m.containerHeight = height
// 	return m
// }

// func (m Model) ContainerWidth() int {
// 	return m.containerWidth
// }

// func (m Model) ContainerHeight() int {
// 	return m.containerHeight
// }
