package introspector

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	events   []string
	viewport viewport.Model
	clientID string
	quit     chan bool
	client   *compiler_introspection.Client
}

type eventMsg string
type tickMsg time.Time

// use this for debugging
func RunIntrospector2(ipcDir string) error {
	_ = make(chan bool)
	client, err := compiler_introspection.NewClient(ipcDir)
	if err != nil {
		return err
	}

	_, err = client.Helo()
	if err != nil {
		return fmt.Errorf("Failed to send HELO request: %w", err)
	}

	for evt := range client.EventChan {
		log.Printf("Got event: %v", evt)

	}

	return nil
}

func RunIntrospector(ipcDir string) error {
	quit := make(chan bool)
	client, err := compiler_introspection.NewClient(ipcDir)
	if err != nil {
		return err
	}

	resp, err := client.Helo()
	if err != nil {
		return fmt.Errorf("Failed to send HELO request: %w", err)
	}

	mainProgram := model{
		viewport: viewport.New(100, 100),
		clientID: resp.ClientID,
		quit:     quit,
		client:   client,
	}
	mainProgram.viewport.YOffset = 1

	if err := tea.NewProgram(mainProgram).Start(); err != nil {
		fmt.Printf("Failed to start TUI: %v\n", err)
		os.Exit(1)
	}

	return nil
}

type TickMsg time.Time
type EventMsg introspection.IntrospectionEvent
type NoEventMsg struct{}

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func pollEvent(events chan introspection.IntrospectionEvent) tea.Cmd {
	return func() tea.Msg {
		event := <-events
		return EventMsg(event)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		pollEvent(m.client.EventChan),
		doTick(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1
		newViewport, _ := m.viewport.Update(msg)
		m.viewport = newViewport
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			close(m.quit)
			return m, tea.Quit
		}

	case EventMsg:
		m.events = append(m.events, fmt.Sprintf("New Event %v", msg))
		return m, pollEvent(m.client.EventChan)

	case TickMsg:
		return m, doTick()
	}

	m.viewport.SetContent(m.eventsToString())
	return m, nil
}

func (m model) View() string {
	var builder strings.Builder
	builder.WriteString("ClientID: ")
	builder.WriteString(m.clientID)
	builder.WriteString("\n")
	for _, e := range m.events {
		builder.WriteString(e)
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m model) eventsToString() string {
	var s string
	for _, e := range m.events {
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Render(e) + "\n"
	}
	return s
}
