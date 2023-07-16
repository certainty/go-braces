package ir

import (
	ast "github.com/certainty/go-braces/internal/compiler/frontend/ast/hl"
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

func (c *SSATransformer) Transform(theAST *ast.Source) (*SSA, error) {
	return nil, nil
}
