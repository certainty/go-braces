package eventlog

import (
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	containerWidth  int
	containerHeight int
	viewport        viewport.Model
	events          []compiler_introspection.CompilerIntrospectionEvent
	isVisible       bool
}

func New() Model {
	return Model{
		events:    []compiler_introspection.CompilerIntrospectionEvent{},
		viewport:  viewport.New(0, 0),
		isVisible: false,
	}
}

func (m Model) IsVisible() bool {
	return m.isVisible
}
