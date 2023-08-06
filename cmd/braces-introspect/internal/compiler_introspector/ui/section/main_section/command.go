package main_section

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
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
