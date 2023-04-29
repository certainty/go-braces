package assembly

import "github.com/certainty/go-braces/internal/isa"

type SourceInfo struct {
	// matches a consecutive slice to instructions to the source it orgiginated from
	InstructionOffset uint64
	Count             uint64
	// TODO add the source information type
	Information interface{}
}

type Chunk struct {
	// code segment
	Code []isa.Instruction
	// constant segment
	Constants ConstantPool
	// debug information
	SourceMap []SourceInfo
}

type AssemblyLibrary struct {
	Chunk Chunk

	// additional properties will follow
	Name    string
	ID      string
	Exports []string
	Imports []string
}

// An assembly module is usually the result of compiling a single file
// Muliple assembly modules can be compined into a single assembly
type AssemblyModule struct {
	Chunk     Chunk
	Libraries []AssemblyLibrary
}

type AssemblyType uint8

const (
	AssemblyTypeExecutable AssemblyType = iota
	AssemblyTypeLibrary
)

type Assembly struct {
	Type    AssemblyType
	Name    string
	Modules []AssemblyModule
}
