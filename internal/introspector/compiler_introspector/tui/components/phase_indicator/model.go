package phase_indicator

import (
	"fmt"

	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	containerWidth  int //noilint:unused,structcheck
	containerHeight int //noilint:unused,structcheck
	phases          []compiler_introspection.CompilationPhase
	currentPhase    int
	finished        bool
	theme           theme.Theme
}

func NewModel(theme theme.Theme) Model {
	phases := []compiler_introspection.CompilationPhase{
		compiler_introspection.CompilationPhaseRead,
		compiler_introspection.CompilationPhaseParse,
		compiler_introspection.CompilationPhaseTypeCheck,
		compiler_introspection.CompilationPhaseOptimize,
		compiler_introspection.CompilationPhaseCodegen,
	}
	return Model{
		containerWidth:  0,
		containerHeight: 0,
		currentPhase:    0,
		phases:          phases,
		theme:           theme,
	}
}

func (m *Model) SetCurrentPhase(phase compiler_introspection.CompilationPhase) {
	for i, p := range m.phases {
		if p == phase {
			m.currentPhase = i
		}
	}
}

func (m *Model) Finish() {
	m.finished = true
}

func (m *Model) Reset() {
	m.currentPhase = 0
	m.finished = false
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	phaseLabelWidth := m.containerWidth / len(m.phases)
	phaseLabels := []string{}
	activeLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Subtext0).Foreground(m.theme.Colors.Base).Bold(true).Width(phaseLabelWidth).Align(lipgloss.Center)
	completedPhaseLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Green).Foreground(m.theme.Colors.Base).Width(phaseLabelWidth).Align(lipgloss.Center)
	upcomingPhaseLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Overlay0).Width(phaseLabelWidth).Align(lipgloss.Center).Faint(true)

	for i, phase := range m.phases {
		if i == m.currentPhase && !m.finished {
			phaseLabels = append(phaseLabels, activeLabelStyle.Render(fmt.Sprintf("\u2022 %s", phase)))
		} else if i < m.currentPhase || m.finished {
			phaseLabels = append(phaseLabels, completedPhaseLabelStyle.Render(fmt.Sprintf("  %s", phase)))
		} else {
			phaseLabels = append(phaseLabels, upcomingPhaseLabelStyle.Render(fmt.Sprintf("  %s", phase)))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, phaseLabels...)
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
