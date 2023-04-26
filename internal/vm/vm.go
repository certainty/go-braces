package vm

import "github.com/certainty/go-braces/internal/compiler"
import "github.com/certainty/go-braces/internal/vm/language/value"

type VmOptions struct {
}

type VM struct {
}

func DefaultOptions() VmOptions {
	return VmOptions{}
}

func NewVM(options VmOptions) VM {
	return VM{}
}

func (vm *VM) Execute(compilationUnit *compiler.CompilationUnit) (value.Value, error) {
	return true, nil
}
