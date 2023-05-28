package phase_indicator

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
)

type Model struct {
	containerWidth  int //noilint:unused,structcheck
	containerHeight int //noilint:unused,structcheck
	phases          []compiler_introspection.CompilationPhase
	currentPhase    int
	finished        bool
	theme           theme.Theme
}

func New(theme theme.Theme) Model {
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
