package frame

import (
	"log"
	"time"

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

func CmdBreakpointContinue(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.BreakpointContinue()
		log.Printf("BreakpointContinue: %v", err)

		if err != nil {
			return common.MsgError{Err: err}
		} else {
			return common.MsgRequestStatus{RequestStatus: common.RequestSent}
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
