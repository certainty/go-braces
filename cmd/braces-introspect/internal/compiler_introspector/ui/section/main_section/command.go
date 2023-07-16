package main_section

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	tea "github.com/charmbracelet/bubbletea"
)

func CmdBreakpointContinue(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		err := client.BreakpointContinue()

		if err != nil {
			return common.MsgError{Err: err}
		} else {
			return common.MsgRequestStatus{RequestStatus: common.RequestSent}
		}
	}
}
