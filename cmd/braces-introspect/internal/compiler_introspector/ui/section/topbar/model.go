package topbar

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
)

type Model struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme

	Title         string
	isConnected   bool
	requestStatus common.RequestStatus
}

func New(theme theme.Theme, title string) Model {
	return Model{
		theme:         theme,
		Title:         title,
		isConnected:   false,
		requestStatus: common.NoRequest,
	}
}
