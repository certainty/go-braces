package ir

import (
	"github.com/certainty/go-braces/internal/compiler/frontend/parser/ast"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type SSA struct {
}

type SSATransformer struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewSSATransformer(instrumentation compiler_introspection.Instrumentation) *SSATransformer {
	return &SSATransformer{instrumentation: instrumentation}
}

func (c *SSATransformer) Transform(theAST *ast.AST) (*SSA, error) {
	return nil, nil
}
