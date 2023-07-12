package compiler_introspection

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type IntrospectionSubject interface{}

type CompilationPhase string

const (
	CompilationPhaseRead        CompilationPhase = "read"
	CompilationPhaseParse       CompilationPhase = "parse"
	CompilationPhaseTypeCheck   CompilationPhase = "typecheck"
	CompilationPhaseLowerToCore CompilationPhase = "lowertocore"
	CompilationPhaseLowerToIR   CompilationPhase = "lowertoir"
	CompilationPhaseOptimize    CompilationPhase = "optimize"
	CompilationPhaseCodegen     CompilationPhase = "codegen"
)

type BreakpointID string

const (
	BPCompilerBeforeLex BreakpointID = "compiler:before:lex"
	BPCompilerAfterLex  BreakpointID = "compiler:after:lex"
	BPReaderParseDatum  BreakpointID = "lexer:parse"
	BPReaderAccepted    BreakpointID = "lexer:accepted"

	BPCompilerBeforeParse BreakpointID = "compiler:before:parse"
	BPCompilerAfterParse  BreakpointID = "compiler:after:parse"

	BPCompilerBeforeCoreCompile BreakpointID = "compiler:before:corecompile"

	BPCompilerBeforeTypeCheck BreakpointID = "compiler:before:typecheck"
	BPCompilerAfterTypeCheck  BreakpointID = "compiler:after:typecheck"
	BPCompilerBeforeOptimize  BreakpointID = "compiler:before:optimize"
	BPCompilerAfterOptimize   BreakpointID = "compiler:after:optimize"
	BPCompilerBeforeCodegen   BreakpointID = "compiler:before:codegen"
	BPCompilerAfterCodegen    BreakpointID = "compiler:after:codegen"
)

// events
type CompilerIntrospectionEvent interface{}

type EventBeginCompileModule struct {
	Origin     location.Origin
	SourceCode string
}

func (e EventBeginCompileModule) String() string {
	return fmt.Sprintf("EventBeginCompileModule{Location: %s, SourceCodeSize: %d}", e.Origin.Name(), len(e.SourceCode))
}

type EventEndCompileModule struct {
	Meta      isa.AssemblyMeta
	Functions []isa.Function
	Closures  []isa.Closure
}

func NewEventEndCompileModule(module isa.AssemblyModule) EventEndCompileModule {
	return EventEndCompileModule{
		Meta: module.Meta,
	}
}

func (e EventEndCompileModule) String() string {
	return fmt.Sprintf("EventEndCompileModule{Meta: %v }", e.Meta)
}

type EventEnterPhase struct {
	Phase CompilationPhase
}

type EventLeavePhase struct {
	Phase CompilationPhase
}

func (e EventEnterPhase) String() string {
	return fmt.Sprintf("EventEnterPhase{Phase: %s}", e.Phase)
}

func (e EventLeavePhase) String() string {
	return fmt.Sprintf("EventLeavePhase{Phase: %s}", e.Phase)
}

type EventBreakpoint struct {
	ID BreakpointID
}

func (e EventBreakpoint) String() string {
	return fmt.Sprintf("EventBreakpoint{ID: %s}", e.ID)
}

// control stuff

type CompilerIntrospectionControl interface{}

// TODO: do we need a correlation ID?
type CommandOk struct {
	Value CompilerIntrospectionControl
}

type CommandError struct {
	Message string
}                                // error
type BreakpointContinue struct{} // continue execution

type Instrumentation interface {
	EnterPhase(phase CompilationPhase)
	LeavePhase(phase CompilationPhase)
	EnterCompilerModule(origin location.Origin, sourceCode string)
	LeaveCompilerModule(module isa.AssemblyModule)
	Breakpoint(id BreakpointID, subject IntrospectionSubject)
}

type InstrumentationFromServer struct {
	haltOnBreakpoint bool
	server           *Server
}

func NewInstrumentationFromServer(server *Server) Instrumentation {
	return &InstrumentationFromServer{server: server, haltOnBreakpoint: true}
}

func (s *InstrumentationFromServer) EnterPhase(phase CompilationPhase) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(EventEnterPhase{Phase: phase}) // nolint:errcheck
	}
}

func (s *InstrumentationFromServer) LeavePhase(phase CompilationPhase) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(EventLeavePhase{Phase: phase}) //nolint:errcheck
	}
}

func (s *InstrumentationFromServer) EnterCompilerModule(origin location.Origin, sourceCode string) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(EventBeginCompileModule{Origin: origin, SourceCode: sourceCode}) //nolint:errcheck
	}
}

func (s *InstrumentationFromServer) LeaveCompilerModule(module isa.AssemblyModule) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(NewEventEndCompileModule(module)) //nolint:errcheck
	}
}

func (s *InstrumentationFromServer) Breakpoint(id BreakpointID, subject IntrospectionSubject) {
	if !s.haltOnBreakpoint {
		return
	}

	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(EventBreakpoint{id}) //nolint:errcheck
		s.breakpointRepl(id, subject)
		log.Printf("Breakpoint resumed  %s", id)
	}
}

func (s *InstrumentationFromServer) breakpointRepl(id BreakpointID, subject IntrospectionSubject) {
	log.Printf("Breakpoint %s hit. Entering REPL for breakpoint", id)
	log.Printf("Subject: %v", subject)

	for {
		nextCommand, err := s.server.ReceiveControl()
		if err != nil {
			log.Printf("Error receiving control: %s", err)
			continue
		}

		switch nextCommand.(type) {
		case BreakpointContinue:
			log.Printf("Continuing execution")
			s.server.SendControl(CommandOk{})
			return
		default:
			log.Printf("Unknown command: %v", nextCommand)
		}
	}
}
