package compiler_introspection

import (
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/shared/isa"
)

type NullInstrumentation struct{}

func NewNullInstrumentation() *NullInstrumentation {
	return &NullInstrumentation{}
}

func (n NullInstrumentation) EnterPhase(phase CompilationPhase) {}

func (n *NullInstrumentation) LeavePhase(phase CompilationPhase) {}

func (s *NullInstrumentation) EnterCompilerModule(origin token.Origin, sourceCode string) {}

func (s *NullInstrumentation) LeaveCompilerModule(module isa.AssemblyModule) {}

func (n NullInstrumentation) Breakpoint(id BreakpointID, subject IntrospectionSubject) {}
