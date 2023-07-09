package isa

type AssemblyModule struct {
	Meta       AssemblyMeta
	Closures   []Closure
	Functions  []Function
	EntryPoint int

	// will be refined later
	Export  []interface{}
	Imports []interface{}
}

func NewAssemblyModule(meta AssemblyMeta, closures []Closure, funcs []Function, exports []interface{}, imports []interface{}) *AssemblyModule {
	return &AssemblyModule{
		Meta:      meta,
		Closures:  closures,
		Functions: funcs,
		// TODO: this is a hack
		EntryPoint: -1,
		Export:     exports,
		Imports:    imports,
	}
}
