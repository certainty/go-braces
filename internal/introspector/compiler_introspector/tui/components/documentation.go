package components

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type DocumentationDialogModel struct {
	ContainerWidth  int
	ContainerHeight int
	theme           theme.Theme
	Active          bool
	Content         string
}

func InitialDocumentationDialogModel(theme theme.Theme) DocumentationDialogModel {
	return DocumentationDialogModel{
		theme:  theme,
		Active: false,
	}
}

func (m DocumentationDialogModel) Init() tea.Cmd {
	return nil
}

func (m DocumentationDialogModel) Update(msg tea.Msg) (DocumentationDialogModel, tea.Cmd) {
	return m, nil
}

func (m DocumentationDialogModel) View() string {
	return ""
}
