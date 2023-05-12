package compiler_introspection

import (
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type NullInstrumentation struct{}

func NewNullInstrumentation() *NullInstrumentation {
	return &NullInstrumentation{}
}

func (n NullInstrumentation) EnterPhase(phase CompilationPhase) error {
	return nil
}

func (n *NullInstrumentation) LeavePhase(phase CompilationPhase) error {
	return nil
}

func (s *NullInstrumentation) EnterCompilerModule(origin location.Origin, sourceCode string) error {
	return nil
}

func (s *NullInstrumentation) LeaveCompilerModule(module isa.AssemblyModule) error {
	return nil
}

func (n NullInstrumentation) Breakpoint(id BreakpointID, subject IntrospectionSubject) error {
	return nil
}
