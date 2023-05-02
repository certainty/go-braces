package isa

type AssemblyModule struct {
	Meta AssemblyMeta
	Code *CodeUnit

	Closures []ClosureValue

	// will be refined later
	Export  []interface{}
	Imports []interface{}
}

func NewAssemblyModule(meta AssemblyMeta, code *CodeUnit, closures []ClosureValue, exports []interface{}, imports []interface{}) *AssemblyModule {
	return &AssemblyModule{
		Meta:     meta,
		Code:     code,
		Closures: closures,
		Export:   exports,
		Imports:  imports,
	}
}
