package isa

type AssemblyModule struct {
	Meta       AssemblyMeta
	Closures   []Closure
	Functions  []Function
	EntryPoint int
}

func NewAssemblyModule(meta AssemblyMeta) *AssemblyModule {
	return &AssemblyModule{
		Meta:       meta,
		Closures:   []Closure{},
		Functions:  []Function{},
		EntryPoint: -1,
	}
}
