package components

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NewEventMessage compiler_introspection.CompilerIntrospectionEvent

type EventLogModel struct {
	ContainerWidth  int
	ContainerHeight int
	Events          []compiler_introspection.CompilerIntrospectionEvent
	viewport        viewport.Model
}

func InitialEventLogModel() EventLogModel {
	return EventLogModel{
		Events:   []compiler_introspection.CompilerIntrospectionEvent{},
		viewport: viewport.New(0, 0),
	}
}

func (m EventLogModel) Init() tea.Cmd {
	return nil
}

func (m EventLogModel) Update(msg tea.Msg) (EventLogModel, tea.Cmd) {
	switch msg := msg.(type) {
	case NewEventMessage:
		m.Events = append(m.Events, msg)
	}
	return m, nil
}

func (m EventLogModel) View() string {
	m.viewport.Width = m.ContainerWidth
	m.viewport.Height = m.ContainerHeight - 2 // header

	eventContent := []string{}
	for _, event := range m.Events {
		nextEvent := fmt.Sprintf("%v", event)
		eventContent = append(eventContent, nextEvent)
	}

	m.viewport.SetContent(strings.Join(eventContent, "\n"))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Bold(true).Render("Events"),
		"",
		m.viewport.View(),
	)
}
