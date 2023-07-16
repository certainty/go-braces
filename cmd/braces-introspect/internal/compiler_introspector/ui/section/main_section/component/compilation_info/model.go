package compilation_info

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/token"
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
)

type Model struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme
	Origin          *token.Origin
	CompilerOptions []string
}

func New(theme theme.Theme, input *input.Input, options []string) Model {
	return Model{
		theme:           theme,
		CompilerOptions: options,
	}
}
