package frame

import (
	"log"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
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
		updatedEventlog, _ := m.eventLog.Update(eventlog.MsgToggleVisibility{})
		m.eventLog = updatedEventlog.(eventlog.Model)
	}
	return m, nil
}

func (m model) propagateResize() tea.Model {
	updatedTopBar, _ := m.topBar.Update(common.MsgResize{Width: m.width, Height: 1})
	m.topBar = updatedTopBar.(topbar.Model)

	updatedEventLog, _ := m.eventLog.Update(common.MsgResize{Width: m.width, Height: 20})
	m.eventLog = updatedEventLog.(eventlog.Model)

	updatedStatus, _ := m.statusBar.Update(common.MsgResize{Width: m.width, Height: 1})
	m.statusBar = updatedStatus.(statusbar.Model)

	return m
}

func (m model) propagateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd     tea.Cmd
		cmds    []tea.Cmd
		updated tea.Model
	)

	isConnected := m.client.IsConnected()

	//topbar
	updated, cmd = m.topBar.Update(msg)
	m.topBar = updated.(topbar.Model)
	cmds = append(cmds, cmd)

	updated, cmd = m.topBar.Update(common.MsgClientConnected(isConnected))
	m.topBar = updated.(topbar.Model)
	cmds = append(cmds, cmd)

	// eventlog
	updated, cmd = m.eventLog.Update(msg)
	m.eventLog = updated.(eventlog.Model)
	cmds = append(cmds, cmd)

	// statusbar
	updated, cmd = m.statusBar.Update(msg)
	m.statusBar = updated.(statusbar.Model)
	cmds = append(cmds, cmd)

	updated, cmd = m.statusBar.Update(common.MsgClientConnected(isConnected))
	m.statusBar = updated.(statusbar.Model)
	cmds = append(cmds, cmd)

	return m, nil
}
