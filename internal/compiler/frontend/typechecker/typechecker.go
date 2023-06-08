package typechecker

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type TypeChecker struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewTypeChecker(Instrumentation compiler_introspection.Instrumentation) *TypeChecker {
	return &TypeChecker{instrumentation: Instrumentation}
}

func (t *TypeChecker) Check(ast *parser.AST) error {
	t.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseTypeCheck)
	defer t.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseTypeCheck)

	return nil
}
