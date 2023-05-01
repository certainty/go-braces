package isa

const CURRENT_ABI_VERSION uint64 = 1

type AssemblyType uint8

const (
	AssemblyTypeExecutable AssemblyType = iota
	AssemblyTypeLibrary
)

type AssemblyMeta struct {
	ABIVersion uint64
	Name       string
	Type       AssemblyType
}

func NewAssemblyMeta(name string, t AssemblyType) AssemblyMeta {
	return AssemblyMeta{
		ABIVersion: CURRENT_ABI_VERSION,
		Name:       name,
		Type:       t,
	}
}
