package vm

import (
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa/assembly"
	"github.com/certainty/go-braces/internal/vm/language/value"
)

type VmOptions struct {
	introspectionAPI *introspection.API
}

type VM struct {
	introspectionAPI introspection.API
}

func DefaultOptions() VmOptions {
	return VmOptions{}
}

func NewVM(options VmOptions) VM {
	if options.introspectionAPI == nil {
		return VM{introspectionAPI: introspection.NullAPI()}
	} else {
		return VM{introspectionAPI: *options.introspectionAPI}
	}
}

func (vm *VM) ExecuteModule(module *assembly.AssemblyModule) (value.Value, error) {
	return true, nil
}
