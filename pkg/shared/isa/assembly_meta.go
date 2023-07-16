package isa

import "fmt"

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

func (a AssemblyMeta) String() string {
	return fmt.Sprintf("AssemblyMeta{ABIVersion: %d, Name: %s, Type: %d}", a.ABIVersion, a.Name, a.Type)
}
