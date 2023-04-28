package codegen

import (
	"github.com/certainty/go-braces/internal/introspection"
)

type Codegenerator struct {
	introspectionAPI *introspection.API
}

func NewCodegenerator(introspectionAPI *introspection.API) *Codegenerator {
	return &Codegenerator{introspectionAPI: introspectionAPI}
}
