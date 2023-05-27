package compile_activity

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/activities"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components/phase_indicator"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	containerHeight int
	containerWidth  int
	theme           theme.Theme
	phase           phase_indicator.Model
}

func NewModel(theme theme.Theme) Model {
	return Model{
		phase: phase_indicator.NewModel(theme),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case messages.IntrospectionEventMsg:
		switch evt := msg.Event.(type) {
		case compiler_introspection.EventBeginCompileModule:
			m.phase.Reset()
		case compiler_introspection.EventEnterPhase:
			m.phase.SetCurrentPhase(evt.Phase)
		case compiler_introspection.EventEndCompileModule:
			m.phase.Finish()
		}
	}
	return m, nil
}

func (m Model) View() string {
	phaseIndicator := m.phase.View()
	mainView := lipgloss.NewStyle().Height(m.containerHeight).Width(m.containerWidth).Render("Compilation Activity")

	return lipgloss.JoinVertical(lipgloss.Top, phaseIndicator, mainView)
}

func (m Model) Resize(width, height int) activities.Model {
	m.phase = m.phase.Resize(width, 3).(phase_indicator.Model)
	m.containerWidth = width
	m.containerHeight = height

	return m
}

func (m Model) ContainerWidth() int {
	return m.containerWidth
}

func (m Model) ContainerHeight() int {
	return m.containerHeight
}
