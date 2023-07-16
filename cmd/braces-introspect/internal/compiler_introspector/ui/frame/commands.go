package frame

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	tea "github.com/charmbracelet/bubbletea"
	"time"
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

func CmdCheckClientConnection(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		connected := client.IsConnected()
		return common.MsgClientConnected(connected)
	}
}
