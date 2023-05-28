package read

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/phase/read/component/source_code"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		return m.propagateResize(), nil
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

func (m Model) propagateResize() tea.Model {
	updatedSourceCodePane, _ := m.sourceCodePane.Update(common.MsgResize{Height: m.containerHeight, Width: m.containerWidth / 3})
	m.sourceCodePane = updatedSourceCodePane.(source_code.Model)

	return m
}
