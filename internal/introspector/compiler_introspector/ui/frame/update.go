package frame

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	propagateUpdates := true

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		propagateUpdates = false
		m.propagateResize()
	case common.MsgTick:
		cmds = append(cmds, common.CmdTick())
	case tea.KeyMsg:
		cmd = m.handleKeys(msg)
		cmds = append(cmds, cmd)
	case common.MsgClientConnected:
		if bool(msg) {
			m.status = Connected
		} else {
			m.status = Disconnected
		}
		cmds = append(cmds, CmdCheckClientConnection(m.client), CmdGetEvent(m.client))
	case common.MsgIntrospectionEvent:
		cmds = append(cmds, CmdGetEvent(m.client))
	case common.MsgError:
		m.status = Error
		m.err = msg.Err
	}

	if propagateUpdates {
		cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return tea.Quit
	case key.Matches(msg, m.keyMap.ToggleEventLog):
		m.updateSection(SectionEventLog, eventlog.MsgToggleVisibility{})
		m.propagateResize()
	case key.Matches(msg, m.keyMap.Continue):
		return CmdBreakpointContinue(m.client)
	default:
		return nil
	}

	return nil
}

func (m *model) propagateResize() {
	m.updateSection(SectionTopBar, common.NewMsgResize(m.width, 1))
	m.updateSection(SectionEventLog, common.NewMsgResize(m.width, 20))
	m.updateSection(SectionStatusBar, common.NewMsgResize(m.width, 1))

	mainContentHeight := m.height - 2 //statusBar, topBar
	if m.sections[SectionEventLog].(eventlog.Model).IsVisible() {
		mainContentHeight -= 20
	}
	m.updateSection(SectionMain, common.NewMsgResize(m.width, mainContentHeight))
}

func (m *model) propagateUpdate(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	for idx := range m.sections {
		cmd = m.updateSection(Section(idx), msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (m *model) updateSection(idx Section, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.sections[idx], cmd = m.sections[idx].Update(msg)
	return cmd
}
