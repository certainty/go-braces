package compilation_info

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/theme"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/lexer"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
)

type Model struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme
	Origin          *token.Origin
	CompilerOptions []string
}

func New(theme theme.Theme, input *lexer.Input, options []string) Model {
	return Model{
		theme:           theme,
		CompilerOptions: options,
	}
}
