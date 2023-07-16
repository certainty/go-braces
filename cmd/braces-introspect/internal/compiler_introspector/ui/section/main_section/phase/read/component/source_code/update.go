package source_code

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/philistino/teacup/code"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.codeViewer, cmd = m.codeViewer.Update(msg)

	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		m.codeViewer.Height = msg.Height - 1
		m.codeViewer.Width = msg.Width - 1
	case MsgUpdateSourceCode:
		m.setSourceCode(string(msg))
	}
	return m, cmd
}

func (m *Model) setSourceCode(sourceCode string) error {
	highlightedCode, err := code.Highlight(sourceCode, "scm", "catppuccin-macchiato")
	m.code = m.prependLineNumbers(highlightedCode)

	if err != nil {
		return err
	}
	return nil
}

func (m Model) prependLineNumbers(code string) string {
	lines := strings.Split(code, "\n")
	var numberedLines []string
	for i, line := range lines {
		number := m.lineNumbersStyle.Render(fmt.Sprintf("%d", i+1))
		numberedLines = append(numberedLines, fmt.Sprintf("%s %s", number, line))
	}
	return strings.Join(numberedLines, "\n")
}
