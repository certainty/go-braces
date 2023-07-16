package eventlog

import (
	"log"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MsgToggleVisibility:
		m.isVisible = !m.isVisible
		log.Printf("m.isVisible = %v", m.isVisible)
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
	case common.MsgIntrospectionEvent:
		m.events = append(m.events, msg)
	}
	return m, nil
}
