package compiler_introspection

import (
	"log"
)

type IntrospectionSubject interface{}
type CompilationPhase string
type EventType string

const (
	EventEnterPhase EventType = "urn:x-braces:compiler:event:enterphase"
	EventLeavePhase EventType = "urn:x-braces:compiler:event:leavephase"
	EventBreakpoint EventType = "urn:x-braces:compiler:event:breakpoint"
)

type CompilerIntrospectionEvent struct {
	Type EventType
	Data interface{}
}

func NewEventEnterPhase(phase CompilationPhase) CompilerIntrospectionEvent {
	return CompilerIntrospectionEvent{
		Type: EventEnterPhase,
		Data: phase,
	}
}

func NewEventLeavePhase(phase CompilationPhase) CompilerIntrospectionEvent {
	return CompilerIntrospectionEvent{
		Type: EventLeavePhase,
		Data: phase,
	}
}

func NewEventBreakpoint(id BreakpointID) CompilerIntrospectionEvent {
	return CompilerIntrospectionEvent{
		Type: EventBreakpoint,
		Data: id,
	}
}

const (
	CompilationPhaseRead      CompilationPhase = "urn:x-braces:compiler:phase:read"
	CompilationPhaseParse     CompilationPhase = "urn:x-braces:compiler:phase:parse"
	CompilationPhaseTypeCheck CompilationPhase = "urn:x-braces:compiler:phase:typecheck"
	CompilationPhaseLowerToIR CompilationPhase = "urn:x-braces:compiler:phase:lowertoir"
	CompilationPhaseSSA       CompilationPhase = "urn:x-braces:compiler:phase:ssa"
	CompilationPhaseOptimize  CompilationPhase = "urn:x-braces:compiler:phase:optimize"
	CompilationPhaseCodegen   CompilationPhase = "urn:x-braces:compiler:phase:codegen"
)

type BreakpointID string

const (
	BPScannerBeginLex BreakpointID = "urn:x-braces:compiler:bp:begin-lex"
	BPScannerEndLex   BreakpointID = "urn:x-braces:compiler:bp:end-lex"
	BPReaderAccepted  BreakpointID = "urn:x-braces:compiler:bp:token-accepted"

	BPCompilerBeforeParse BreakpointID = "urn:x-braces:compiler:bp:before-parse"
	BPCompilerAfterParse  BreakpointID = "urn:x-braces:compiler:bp:after-parse"
)

type CompilerIntrospectionControl interface{}

type CommandOk struct {
	Value CompilerIntrospectionControl
}

type CommandError struct {
	Message string
} // error

type BreakpointContinue struct{} // continue execution

type Instrumentation interface {
	EnterPhase(phase CompilationPhase)
	LeavePhase(phase CompilationPhase)
	Breakpoint(id BreakpointID, subject IntrospectionSubject)
}

type InstrumentationFromServer struct {
	haltOnBreakpoint bool
	server           *Server
}

func NewInstrumentationFromServer(server *Server) Instrumentation {
	return &InstrumentationFromServer{
		server:           server,
		haltOnBreakpoint: true,
	}
}

func (s *InstrumentationFromServer) EnterPhase(phase CompilationPhase) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(NewEventEnterPhase(phase)) // nolint:errcheck
	}
}

func (s *InstrumentationFromServer) LeavePhase(phase CompilationPhase) {
	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(NewEventLeavePhase(phase)) //nolint:errcheck
	}
}

func (s *InstrumentationFromServer) Breakpoint(id BreakpointID, subject IntrospectionSubject) {
	if !s.haltOnBreakpoint {
		return
	}

	if s.server != nil && s.server.IsConnected() {
		s.server.SendEvents(NewEventBreakpoint(id)) //nolint:errcheck
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
