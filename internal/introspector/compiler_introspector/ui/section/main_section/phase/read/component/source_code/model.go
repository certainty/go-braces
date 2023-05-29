package source_code

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	containerWidth, containerHeight int
	theme                           theme.Theme
	codeViewer                      viewport.Model
	codeContentStyle                lipgloss.Style
	lineNumbersStyle                lipgloss.Style
	code                            string
}

func New(theme theme.Theme) Model {
	codeContentStyle := lipgloss.NewStyle().PaddingLeft(1).Border(lipgloss.NormalBorder())
	lineNumbersStyle := lipgloss.NewStyle().Foreground(theme.Colors.Overlay0)

	return Model{
		containerWidth:   0,
		containerHeight:  0,
		theme:            theme,
		code:             "",
		codeViewer:       viewport.New(0, 0),
		codeContentStyle: codeContentStyle,
		lineNumbersStyle: lineNumbersStyle,
	}
}
