package main_section

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/compilation_info"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
)

type Model struct {
	containerWidth, containerHeight int
	theme                           theme.Theme
	compilationInfo                 compilation_info.Model
	phaseIndicator                  phase_indicator.Model
	isCompiling                     bool
}

func New(theme theme.Theme) Model {
	return Model{
		compilationInfo: compilation_info.New(theme, nil, nil),
		phaseIndicator:  phase_indicator.New(theme),
		isCompiling:     false,
	}
}
