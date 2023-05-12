package components

import (
	"strings"

	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InfoModel struct {
	ContainerWidth  int
	ContainerHeight int
	theme           theme.Theme
	Origin          *location.Origin
	CompilerOptions []string
}

func InitialInfoModel(theme theme.Theme, input *input.Input, options []string) InfoModel {
	return InfoModel{
		theme:           theme,
		CompilerOptions: options,
	}
}

func (m InfoModel) Init() tea.Cmd {
	return nil
}

func (m InfoModel) Update(msg tea.Msg) (InfoModel, tea.Cmd) {
	switch msg := msg.(type) {
	case messages.IntrospectionEventMsg:
		switch msg := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			m.Origin = &msg.Origin
		default:
			break
		}
	default:
		break
	}

	return m, nil
}

func (m InfoModel) View() string {
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

	width := m.ContainerWidth
	height := m.ContainerHeight - 2

	return lipgloss.
		NewStyle().
		Border(lipgloss.RoundedBorder(), true, false, true, false).
		BorderForeground(m.theme.Colors.InactiveBorder).
		Width(width).
		Height(height).
		Render(content)
}
