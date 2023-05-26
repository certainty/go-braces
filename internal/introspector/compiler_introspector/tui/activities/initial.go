package activities

import (
	"log"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type InitialActivity struct {
	model components.MasterDetailModel
}

func NewInitialActivity(theme theme.Theme) *InitialActivity {
	return &InitialActivity{
		model: components.InitialMainModel(theme),
	}
}

func (a InitialActivity) Name() string {
	return "Initial"
}

func (a InitialActivity) Model() tea.Model {
	return a.model
}

func (a *InitialActivity) UpdateSize(width, height int) {
	a.model.ContainerWidth = width
	a.model.ContainerHeight = height

	log.Printf("InitialActivity.UpdateSize: width=%d, height=%d\n", a.model.ContainerWidth, a.model.ContainerHeight)
}
