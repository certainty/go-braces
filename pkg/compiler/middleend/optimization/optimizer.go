package optimization

import (
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ssa"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
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

	ssaTransformer := ssa.NewTransformer(o.instrumentation)
	err := ssaTransformer.Transform(intermediate)

	if err != nil {
		return nil, err
	}

	return intermediate, nil
}
