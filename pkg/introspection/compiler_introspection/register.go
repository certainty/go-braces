package compiler_introspection

import (
	"encoding/gob"

	"github.com/certainty/go-braces/internal/compiler/frontend/token"
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
	gob.Register(EventBreakpoint{})

	// events

	// contol stuff
	gob.Register(BreakpointContinue{})
	gob.Register(CommandOk{})
	gob.Register(CommandError{})

	// data
	gob.Register(token.StringOrigin{})
	gob.Register(token.FileOrigin{})
	gob.Register(token.ReplOrigin{})
	gob.Register(isa.AssemblyMeta{})
	gob.Register(isa.CodeUnit{})
	gob.Register(isa.Instruction{})
	gob.Register(isa.OpCode(0))
	gob.Register(isa.Register(0))
	gob.Register(isa.Function{})
	gob.Register(isa.Closure{})
}
