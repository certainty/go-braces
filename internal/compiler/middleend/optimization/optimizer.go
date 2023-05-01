package optimization

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/introspection"
)

type Optimizer struct {
	introspectionAPI introspection.API
}

func NewOptimizer(introspectionAPI introspection.API) *Optimizer {
	return &Optimizer{
		introspectionAPI: introspectionAPI,
	}
}

func (o *Optimizer) Optimize(intermediate *ir.IR) (*ir.IR, error) {
	return nil, nil
}
