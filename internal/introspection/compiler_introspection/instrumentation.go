package compiler_introspection

import (
	"fmt"

	"github.com/certainty/go-braces/internal/compiler/location"
	"github.com/certainty/go-braces/internal/isa"
)

type IntrospectionSubject interface{}

type CompilationPhase string

const (
	CompilationPhaseRead      CompilationPhase = "read"
	CompilationPhaseParse     CompilationPhase = "parse"
	CompilationPhaseTypeCheck CompilationPhase = "typecheck"
	CompilationPhaseOptimize  CompilationPhase = "optimize"
	CompilationPhaseCodegen   CompilationPhase = "codegen"
)

type BreakpointID string

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
	Meta isa.AssemblyMeta
	Code isa.CodeUnit
}

func NewEventEndCompileModule(module isa.AssemblyModule) EventEndCompileModule {
	return EventEndCompileModule{
		Meta: module.Meta,
		Code: *module.Code,
	}
}

func (e EventEndCompileModule) String() string {
	return fmt.Sprintf("EventEndCompileModule{Meta: %v CodeSize: %d ConstandPoolSize: %d}", e.Meta, len(e.Code.Instructions), len(e.Code.Constants))
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

// control stuff
type Instrumentation interface {
	EnterPhase(phase CompilationPhase) error
	LeavePhase(phase CompilationPhase) error
	EnterCompilerModule(origin location.Origin, sourceCode string) error
	LeaveCompilerModule(module isa.AssemblyModule) error
	Breakpoint(id BreakpointID, subject IntrospectionSubject) error
}

type InstrumentationFromServer struct {
	server *Server
}

func NewInstrumentationFromServer(server *Server) Instrumentation {
	return &InstrumentationFromServer{server}
}

func (s *InstrumentationFromServer) EnterPhase(phase CompilationPhase) error {
	if s.server != nil && s.server.IsConnected() {
		return s.server.SendEvents(EventEnterPhase{Phase: phase})
	}
	return nil
}

func (s *InstrumentationFromServer) LeavePhase(phase CompilationPhase) error {
	if s.server != nil && s.server.IsConnected() {
		return s.server.SendEvents(EventLeavePhase{Phase: phase})
	}

	return nil
}

func (s *InstrumentationFromServer) EnterCompilerModule(origin location.Origin, sourceCode string) error {
	if s.server != nil && s.server.IsConnected() {
		return s.server.SendEvents(EventBeginCompileModule{Origin: origin, SourceCode: sourceCode})
	}
	return nil
}

func (s *InstrumentationFromServer) LeaveCompilerModule(module isa.AssemblyModule) error {
	if s.server != nil && s.server.IsConnected() {
		return s.server.SendEvents(NewEventEndCompileModule(module))
	}
	return nil
}

func (s *InstrumentationFromServer) Breakpoint(id BreakpointID, subject IntrospectionSubject) error {
	return nil
}
