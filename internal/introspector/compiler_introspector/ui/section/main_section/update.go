package main_section

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/compilation_info"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		return m.propagateResize(), nil

	case common.MsgIntrospectionEvent:
		var (
			updatedPhaseIndicator tea.Model = m.phaseIndicator
			updatedInfo           tea.Model = m.compilationInfo
			cmd                   tea.Cmd
			cmds                  []tea.Cmd
		)

		m.phasePanes[m.activePhaseIndex], cmd = m.phasePanes[m.activePhaseIndex].Update(msg)
		cmds = append(cmds, cmd)

		switch evt := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			m.isCompiling = true
			m.isFinished = false
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgReset{})
			cmds = append(cmds, cmd)

			updatedInfo, cmd = m.compilationInfo.Update(compilation_info.MsgNewCompilation{
				Options: []string{},
				Origin:  evt.Origin,
			})
			cmds = append(cmds, cmd)

		case compiler_introspection.EventEndCompileModule:
			m.isCompiling = false
			m.isFinished = true
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgFinish{})
			cmds = append(cmds, cmd)

		case compiler_introspection.EventEnterPhase:
			updatedPhaseIndicator, cmd = m.phaseIndicator.Update(phase_indicator.MsgPhase(evt.Phase))
			cmds = append(cmds, cmd)
		}
		m.phaseIndicator = updatedPhaseIndicator.(phase_indicator.Model)
		m.compilationInfo = updatedInfo.(compilation_info.Model)

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) propagateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedInfo, cmd := m.compilationInfo.Update(msg)
	m.compilationInfo = updatedInfo.(compilation_info.Model)

	return m, cmd
}

func (m Model) propagateResize() tea.Model {
	updatedPhaseIndicator, _ := m.phaseIndicator.Update(common.MsgResize{Width: m.containerWidth, Height: 1})
	m.phaseIndicator = updatedPhaseIndicator.(phase_indicator.Model)

	updatedInfo, _ := m.compilationInfo.Update(common.MsgResize{Width: m.containerWidth, Height: 3})
	m.compilationInfo = updatedInfo.(compilation_info.Model)

	updatedPhasePane, _ := m.phasePanes[m.activePhaseIndex].Update(common.MsgResize{Width: m.containerWidth, Height: m.containerHeight - 4})
	m.phasePanes[m.activePhaseIndex] = updatedPhasePane

	return m
}
