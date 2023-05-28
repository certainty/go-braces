package eventlog

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.isVisible {
		log.Printf("not visible")
		return ""
	}

	m.viewport.Width = m.containerWidth
	m.viewport.Height = m.containerHeight - 3 // header

	eventContent := []string{}
	for _, event := range m.events {
		nextEvent := fmt.Sprintf("%v", event)
		eventContent = append(eventContent, nextEvent)
	}

	m.viewport.SetContent(strings.Join(eventContent, "\n"))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Bold(true).Width(m.containerWidth-2).Render("Events"),
		m.viewport.View(),
	)
}
