package isa

type AssemblyModule struct {
	Meta AssemblyMeta
	Code CodeUnit

	Closures []ClosureValue

	// will be refined later
	Export  []interface{}
	Imports []interface{}
}
