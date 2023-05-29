package frame

import (
	"log"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/statusbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/topbar"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds    []tea.Cmd
		cmd     tea.Cmd
		updated tea.Model
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		updated = m.propagateResize()
	case common.MsgTick:
		cmds = append(cmds, common.CmdTick())
		updated, cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		updated, cmd = m.handleKeys(msg)
		cmds = append(cmds, cmd)
	case common.MsgClientConnected:
		log.Printf("Client connected polling for event")
		if msg {
			m.status = Connected
		} else {
			m.status = Disconnected
		}
		cmds = append(cmds, CmdGetEvent(m.client))
		updated, cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	case common.MsgIntrospectionEvent:
		cmds = append(cmds, CmdGetEvent(m.client))
		updated, cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	case common.MsgError:
		m.status = Error
		m.err = msg.Err
		updated, cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	default:
		updated, cmd = m.propagateUpdate(msg)
		cmds = append(cmds, cmd)
	}

	return updated, tea.Batch(cmds...)
}

func (m model) handleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keyMap.ToggleEventLog):
		updatedEventlog, _ := m.sectionEventLog.Update(eventlog.MsgToggleVisibility{})
		m.sectionEventLog = updatedEventlog.(eventlog.Model)
		m = m.propagateResize().(model)
		// TODO: handle this in the main section
		// use an event to load the available shortcuts
	case key.Matches(msg, m.keyMap.Continue):
		return m, CmdBreakpointContinue(m.client)
	}
	return m, nil
}

func (m model) propagateResize() tea.Model {
	updatedTopBar, _ := m.sectionTopBar.Update(common.MsgResize{Width: m.width, Height: 1})
	m.sectionTopBar = updatedTopBar.(topbar.Model)

	updatedEventLog, _ := m.sectionEventLog.Update(common.MsgResize{Width: m.width, Height: 20})
	m.sectionEventLog = updatedEventLog.(eventlog.Model)

	updatedStatus, _ := m.sectionStatusBar.Update(common.MsgResize{Width: m.width, Height: 1})
	m.sectionStatusBar = updatedStatus.(statusbar.Model)

	mainContentHeight := m.height - 2 //statusBar, topBar
	if m.sectionEventLog.IsVisible() {
		mainContentHeight -= 20
	}

	updatedMainContent, _ := m.sectionMain.Update(common.MsgResize{Width: m.width, Height: mainContentHeight})
	m.sectionMain = updatedMainContent.(main_section.Model)

	return m
}

func (m model) propagateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	isConnected := m.client.IsConnected()

	for idx, section := range m.sections {
		m.sections[idx], cmd = section.Update(common.MsgClientConnected(isConnected))
		cmds = append(cmds, cmd)

		m.sections[idx], cmd = section.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, nil
}
