package main_section

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	if !m.isCompiling && !m.isFinished {
		return m.waitingView()
	}

	phaseIndicatorView := m.phaseIndicator.View()
	infoView := m.compilationInfo.View()
	phaseView := m.phasePanes[m.activePhaseIndex].View()

	return lipgloss.JoinVertical(lipgloss.Top, infoView, phaseIndicatorView, phaseView)
}

func (m Model) waitingView() string {
	style := lipgloss.NewStyle().
		Height(m.containerHeight).
		Width(m.containerWidth).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)

	return style.Render("Waiting for next compilation cycle. As soon as the connected compiler begins, the view will change to the corresponding phase introspection.")
}
