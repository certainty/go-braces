package commands

import (
	"time"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
	tea "github.com/charmbracelet/bubbletea"
)

func DoTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return messages.TickMsg(t)
	})
}

func DoGetEvent(client *compiler_introspection.Client) tea.Cmd {
	return func() tea.Msg {
		nextEvent, err := client.PollEvent()

		if err != nil {
			return nil // TODO: provider error message
		} else {
			return messages.IntrospectionEventMsg{Event: nextEvent}
		}
	}
}