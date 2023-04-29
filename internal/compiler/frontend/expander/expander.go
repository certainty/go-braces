package expander

import "github.com/certainty/go-braces/internal/introspection"

type Expander struct {
	introspectionAPI introspection.API
}

func NewExpander(introspectionAPI introspection.API) *Expander {
	return &Expander{introspectionAPI: introspectionAPI}
}
