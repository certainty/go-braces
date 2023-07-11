package isa

type Label string
type Address uint64

type Assembly interface {
	Meta() AssemblyMeta
}

type ExecutableAssembly struct {
	meta       AssemblyMeta
	EntryPoint Address
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

func NewExecutableAssembly(name string, entrypoint Address, modules []AssemblyModule) ExecutableAssembly {
	return ExecutableAssembly{
		meta:       NewAssemblyMeta(name, AssemblyTypeExecutable),
		EntryPoint: Address(entrypoint),
		Modules:    modules,
	}
}

func NewLibraryAssembly(name string, modules []AssemblyModule) LibraryAssembly {
	return LibraryAssembly{
		meta:    NewAssemblyMeta(name, AssemblyTypeLibrary),
		Modules: modules,
	}
}
