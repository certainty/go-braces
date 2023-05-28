package main_section

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/compilation_info"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/phase/read"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	containerWidth, containerHeight int
	theme                           theme.Theme
	compilationInfo                 compilation_info.Model
	phaseIndicator                  phase_indicator.Model
	phasePanes                      []tea.Model
	activePhaseIndex                int
	isCompiling                     bool
	isFinished                      bool
}

const (
	readPhaseIndex      = 0
	compilePhaseIndex   = 1
	typeCheckPhaseIndex = 2
	optimizePhaseIndex  = 3
	codeGenPhaseIndex   = 4
)

func New(theme theme.Theme) Model {
	phasePanes := []tea.Model{
		read.New(theme),
	}

	return Model{
		compilationInfo:  compilation_info.New(theme, nil, nil),
		isCompiling:      false,
		phaseIndicator:   phase_indicator.New(theme),
		isFinished:       false,
		phasePanes:       phasePanes,
		activePhaseIndex: 0,
	}
}
