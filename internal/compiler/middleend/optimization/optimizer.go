package optimization

import (
	"log"

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

func (o *Optimizer) Optimize(intermediate *ir.IR) (*ir.IR, error) {
	o.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseOptimize)
	defer o.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseOptimize)

	log.Printf("opitmizing %v", intermediate)
	return intermediate, nil
}
