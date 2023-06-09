package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

func CmdTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return MsgTick(t)
	})
}

func CmdActivateKeyMap(keyMap KeyMap) tea.Cmd {
	return func() tea.Msg {
		return MsgActivateKeyMap(keyMap)
	}
}
