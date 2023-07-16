package read

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/main_section/phase/read/component/source_code"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/theme"
)

type Model struct {
	containerWidth, containerHeight int
	theme                           theme.Theme
	sourceCodePane                  source_code.Model
}

func New(theme theme.Theme) Model {
	return Model{
		containerWidth:  0,
		containerHeight: 0,
		theme:           theme,
		sourceCodePane:  source_code.New(theme),
	}
}
