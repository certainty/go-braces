package optimization

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/ir"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type Optimizer struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewOptimizer(instrumentation compiler_introspection.Instrumentation) *Optimizer {
	return &Optimizer{
		instrumentation: instrumentation,
	}
}

func (o *Optimizer) Optimize(intermediate *ir.Module) (*ir.Module, error) {
	o.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseOptimize)
	defer o.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseOptimize)
	return intermediate, nil
}
