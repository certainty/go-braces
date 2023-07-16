package ir

import (
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
)

type SSA struct {
}

type SSATransformer struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewSSATransformer(instrumentation compiler_introspection.Instrumentation) *SSATransformer {
	return &SSATransformer{instrumentation: instrumentation}
}
