package eventlog

import (
	"fmt"
	"strings"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NewEventMessage compiler_introspection.CompilerIntrospectionEvent

type Model struct {
	containerWidth  int
	containerHeight int
	Events          []compiler_introspection.CompilerIntrospectionEvent
	viewport        viewport.Model
}

func NewModel() Model {
	return Model{
		Events:   []compiler_introspection.CompilerIntrospectionEvent{},
		viewport: viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewEventMessage:
		m.Events = append(m.Events, msg)
	}
	return m, nil
}

func (m Model) View() string {
	m.viewport.Width = m.containerWidth
	m.viewport.Height = m.containerHeight - 2 // header

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

func (m Model) Resize(width, height int) components.Model {
	m.containerWidth = width
	m.containerHeight = height
	return m
}

func (m Model) ContainerHeight() int {
	return m.containerHeight
}

func (m Model) ContainerWidth() int {
	return m.containerWidth
}
