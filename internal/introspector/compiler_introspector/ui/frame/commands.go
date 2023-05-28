package frame

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func CmdGetEvent(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		nextEvent, err := client.PollEvent()

		if err != nil {
			return common.MsgError{Err: err}
		} else {
			return common.MsgIntrospectionEvent{Event: nextEvent}
		}
	}
}

func CmdConnectClient(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.Connect()

		if err != nil {
			return common.MsgError{Err: err}
		} else {
			return common.MsgClientConnected(true)
		}
	}
}
