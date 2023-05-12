package expander

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type Expander struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewExpander(instrumentation compiler_introspection.Instrumentation) *Expander {
	return &Expander{instrumentation: instrumentation}
}

func (e *Expander) Expand(data isa.Datum) (isa.Datum, error) {
	return data, nil
}
