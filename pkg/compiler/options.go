package compiler

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
)

type CompilerOptions struct {
	instrumentation compiler_introspection.Instrumentation
}

func NewCompilerOptions(instrumentation compiler_introspection.Instrumentation) CompilerOptions {
	return CompilerOptions{
		instrumentation: instrumentation,
	}
}

func DefaultOptions() CompilerOptions {
	return CompilerOptions{
		instrumentation: compiler_introspection.NewNullInstrumentation(),
	}
}
