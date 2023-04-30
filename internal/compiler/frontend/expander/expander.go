package expander

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/reader"
	"github.com/certainty/go-braces/internal/introspection"
)

type Expander struct {
	introspectionAPI introspection.API
}

func NewExpander(introspectionAPI introspection.API) *Expander {
	return &Expander{introspectionAPI: introspectionAPI}
}

func (e *Expander) Expand(data reader.Datum) (reader.Datum, error) {
	return nil, nil
}
