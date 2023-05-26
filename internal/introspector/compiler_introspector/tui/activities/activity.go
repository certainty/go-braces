package activities

import tea "github.com/charmbracelet/bubbletea"

type Activity interface {
	Name() string
	Model() tea.Model
	UpdateSize(width, height int)
}
