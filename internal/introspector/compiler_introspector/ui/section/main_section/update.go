package main_section

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/compilation_info"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.MsgResize:
		m.containerHeight = msg.Height
		m.containerWidth = msg.Width
		m.propagateResize()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, KeyContinue):
			return m, CmdBreakpointContinue(m.client)
		}

	case common.MsgIntrospectionEvent:
		var (
			cmd  tea.Cmd
			cmds []tea.Cmd
		)

		m.phasePanes[m.activePhasePane], cmd = m.phasePanes[m.activePhasePane].Update(msg)
		cmds = append(cmds, cmd)

		switch evt := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			cmd = m.onBeginCompileModule(evt)
			cmds = append(cmds, cmd)

		case compiler_introspection.EventEndCompileModule:
			cmd = m.onEndCompileModule(evt)
			cmds = append(cmds, cmd)

		case compiler_introspection.EventEnterPhase:
			cmd = m.onEnterPhase(evt)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *Model) propagateResize() {
	m.updateSection(SectionPhaseIndicator, common.NewMsgResize(m.containerWidth, 1))
	m.updateSection(SectionCompilationInfo, common.NewMsgResize(m.containerWidth, 3))
	m.updatedPhasePane(m.activePhasePane, common.NewMsgResize(m.containerWidth, m.containerHeight-4))
}

func (m *Model) propagateUpdate(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for idx := range m.sections {
		cmd = m.updateSection(section(idx), msg)
		cmds = append(cmds, cmd)
	}

	cmd = m.updatedPhasePane(m.activePhasePane, msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *Model) updateSection(sect section, msg tea.Msg) tea.Cmd {
	updatedSection, cmd := m.sections[sect].Update(msg)
	m.sections[sect] = updatedSection
	return cmd
}

func (m *Model) updatedPhasePane(pane Pane, msg tea.Msg) tea.Cmd {
	updatedPane, cmd := m.phasePanes[pane].Update(msg)
	m.phasePanes[pane] = updatedPane
	return cmd
}

func (m *Model) onBeginCompileModule(evt compiler_introspection.EventBeginCompileModule) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.isCompiling = true
	m.isFinished = false

	// activate keymap for this mode
	cmds = append(cmds, common.CmdActivateKeyMap(CompileKeyMap))

	cmd = m.updateSection(SectionPhaseIndicator, phase_indicator.MsgReset{})
	cmds = append(cmds, cmd)

	cmd = m.updateSection(SectionCompilationInfo, compilation_info.MsgNewCompilation{
		Options: []string{},
		Origin:  evt.Origin,
	})

	return tea.Batch(append(cmds, cmd)...)
}

func (m *Model) onEndCompileModule(evt compiler_introspection.EventEndCompileModule) tea.Cmd {
	m.isCompiling = false
	m.isFinished = true

	return m.updateSection(SectionPhaseIndicator, phase_indicator.MsgFinish{})
}

func (m *Model) onEnterPhase(evt compiler_introspection.EventEnterPhase) tea.Cmd {
	return m.updateSection(SectionPhaseIndicator, phase_indicator.MsgPhase(evt.Phase))
}
