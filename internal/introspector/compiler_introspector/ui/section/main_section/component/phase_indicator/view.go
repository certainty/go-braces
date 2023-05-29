package phase_indicator

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	phaseLabelWidth := m.containerWidth / len(m.phases)
	phaseLabels := []string{}
	activeLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Subtext0).Foreground(m.theme.Colors.Base).Bold(true).Width(phaseLabelWidth).Align(lipgloss.Center).MarginRight(1)
	completedPhaseLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Green).Foreground(m.theme.Colors.Base).Width(phaseLabelWidth).Align(lipgloss.Center).MarginRight(1)
	upcomingPhaseLabelStyle := lipgloss.NewStyle().Background(m.theme.Colors.Overlay0).Width(phaseLabelWidth).Align(lipgloss.Center).Faint(true).MarginRight(1)

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
