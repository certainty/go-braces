package compiler_introspection

import (
	"github.com/certainty/go-braces/internal/compiler"
	"github.com/certainty/go-braces/internal/introspection"
)

type CompilerIntrospectionAPI struct {
	inSingleStep bool
	server       *Server
}

func (c *CompilerIntrospectionAPI) SendEvent(event introspection.IntrospectionEvent) {
	c.server.EventChan <- event
}

func (c *CompilerIntrospectionAPI) SingleStepBarrier(subject introspection.IntrospectionSubject) {
	if !c.inSingleStep {
		return
	}

	compiler, ok := subject.(compiler.Compiler)
	if !ok {
		panic("CompilerIntrospectionAPI.SingleStepBarrier: subject is not a compiler")
	}

	// check type of subject is compiler
	c.server.EventChan <- introspection.EventSingleStepBarrierReached{}
	c.singleStepControlLoop(&compiler)
}

func (c *CompilerIntrospectionAPI) singleStepControlLoop(compiler *compiler.Compiler) {
	// do stuff
}
