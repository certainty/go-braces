package components

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpDialogModel struct {
	ContainerWidth  int
	ContainerHeight int
	theme           theme.Theme
	Active          bool
	Content         string
}

func InitialHelpDialogModel(theme theme.Theme) HelpDialogModel {
	return HelpDialogModel{
		theme:  theme,
		Active: false,
	}
}

func (m HelpDialogModel) Init() tea.Cmd {
	return nil
}

func (m HelpDialogModel) Update(msg tea.Msg) (HelpDialogModel, tea.Cmd) {
	return m, nil
}

func (m HelpDialogModel) View() string {
	return ""
}
