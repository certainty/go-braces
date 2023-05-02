package isa

type Assembly interface {
	Meta() AssemblyMeta
}

type ExecutableAssembly struct {
	meta       AssemblyMeta
	EntryPoint uint64
	Modules    []AssemblyModule
}

type LibraryAssembly struct {
	meta    AssemblyMeta
	Modules []AssemblyModule
}

func (a ExecutableAssembly) Meta() AssemblyMeta {
	return a.meta
}

func (a LibraryAssembly) Meta() AssemblyMeta {
	return a.meta
}

func NewExecutableAssembly(name string, entrypoint uint64, modules []AssemblyModule) ExecutableAssembly {
	return ExecutableAssembly{
		meta:       NewAssemblyMeta(name, AssemblyTypeExecutable),
		EntryPoint: entrypoint,
		Modules:    modules,
	}
}

func NewLibraryAssembly(name string, modules []AssemblyModule) LibraryAssembly {
	return LibraryAssembly{
		meta:    NewAssemblyMeta(name, AssemblyTypeLibrary),
		Modules: modules,
	}
}
