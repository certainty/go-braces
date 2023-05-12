package compiler_introspection

import (
	"encoding/gob"
	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
)

func RegisterTypes() {
	introspection.RegisterWireTypes()
	gob.Register(CompilationPhase(""))
	gob.Register(BreakpointID(""))

	gob.Register(EventEnterPhase{})
	gob.Register(EventLeavePhase{})
	gob.Register(EventBeginCompileModule{})
	gob.Register(EventEndCompileModule{})

	// events

	// contol stuff

	// data
	gob.Register(location.StringOrigin{})
	gob.Register(location.FileOrigin{})
	gob.Register(location.ReplOrigin{})
	gob.Register(isa.AssemblyMeta{})
	gob.Register(isa.CodeUnit{})
	gob.Register(isa.Instruction{})
	gob.Register(isa.OpCode(0))
	gob.Register(isa.Register(0))
	gob.Register(isa.BoolValue(true))
	gob.Register(isa.ProcedureValue{})
	gob.Register(isa.ClosureValue{})
}
