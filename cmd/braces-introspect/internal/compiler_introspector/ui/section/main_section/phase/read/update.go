package read

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		m.propagateResize()
		return m, nil

	case common.MsgIntrospectionEvent:
		switch evt := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			updatedSourceCodePane, cmd := m.sourceCodePane.Update(source_code.MsgUpdateSourceCode(evt.SourceCode))
			m.sourceCodePane = updatedSourceCodePane.(source_code.Model)
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) propagateResize() {
	paneWidth := m.containerWidth / 3

	updatedSourceCodePane, _ := m.sourceCodePane.Update(common.MsgResize{Height: m.containerHeight, Width: paneWidth})
	m.sourceCodePane = updatedSourceCodePane.(source_code.Model)
}
