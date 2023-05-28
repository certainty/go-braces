package main_section

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgIntrospectionEvent:
		var (
			updatedPhaseIndicator tea.Model
			cmd                   tea.Cmd
		)
		switch evt := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			m.isCompiling = true
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgReset{})
		case compiler_introspection.EventEndCompileModule:
			m.isCompiling = false
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgFinish{})
		case compiler_introspection.EventEnterPhase:
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgPhase(evt.Phase))
		default:
			updatedPhaseIndicator = m.phaseIndicator
		}
		m.phaseIndicator = updatedPhaseIndicator.(phase_indicator.Model)
		return m, cmd
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		return m.propagateResize(), nil
	}

	return m, nil
}

func (m Model) propagateResize() tea.Model {
	updatedPhaseIndicator, _ := m.phaseIndicator.Update(common.MsgResize{Width: m.containerWidth, Height: 1})
	m.phaseIndicator = updatedPhaseIndicator.(phase_indicator.Model)

	return m
}
