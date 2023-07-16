package main_section

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/compilation_info"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/component/phase_indicator"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section/phase/read"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	tea "github.com/charmbracelet/bubbletea"
)

type section int

const (
	SectionCompilationInfo section = iota
	SectionPhaseIndicator
)

type Pane int

const (
	PaneReadPhase      = 0
	PaneCompilePhase   = 1
	PaneTypeCheckPhase = 2
	PaneOptimizePhase  = 3
	PaneCodeGenPhase   = 4
)

type Model struct {
	containerWidth, containerHeight int
	// nolint:unused
	theme      theme.Theme
	sections   []tea.Model
	phasePanes []tea.Model
	// nolint:unused
	keyMap common.KeyMap

	client          *compiler_introspection.Client
	activePhasePane Pane
	isCompiling     bool
	isFinished      bool
}

func New(theme theme.Theme, client *compiler_introspection.Client) Model {
	phasePanes := []tea.Model{
		read.New(theme),
	}

	sections := make([]tea.Model, 2)
	sections[SectionCompilationInfo] = compilation_info.New(theme, nil, nil)
	sections[SectionPhaseIndicator] = phase_indicator.New(theme)

	return Model{
		sections:        sections,
		phasePanes:      phasePanes,
		isCompiling:     false,
		isFinished:      false,
		activePhasePane: 0,
		client:          client,
	}
}
