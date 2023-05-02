package vm

import (
	"github.com/certainty/go-braces/internal/introspection"
	"github.com/certainty/go-braces/internal/isa"
)

const (
	REG_GP_COUNT = 256
	REG_SP_COUNT = 16
	REG_SP_HALT  = 0
	REG_SP_ACCU  = 1
)

type VmOptions struct {
	introspectionAPI *introspection.API
}

type VM struct {
	introspectionAPI introspection.API
	registers        [REG_SP_COUNT + REG_GP_COUNT]isa.Value
	code             *isa.CodeUnit
	pc               int
}

func DefaultOptions() VmOptions {
	return VmOptions{}
}

func NewVM(options VmOptions) *VM {
	vm := VM{
		pc: 0,
	}

	if options.introspectionAPI == nil {
		vm.introspectionAPI = introspection.NullAPI()
	} else {
		vm.introspectionAPI = *options.introspectionAPI
	}

	return &vm
}

func (vm *VM) ExecuteModule(module *isa.AssemblyModule) (*isa.Value, error) {
	vm.code = module.Code

	for vm.pc < len((*vm.code).Instructions) {
		instr := (*vm.code).Instructions[vm.pc]
		vm.pc++

		switch instr.Opcode {
		case isa.OP_TRUE:
			vm.registers[REG_SP_ACCU] = isa.BoolValue(true)
		case isa.OP_HALT:
			vm.registers[REG_SP_HALT] = vm.registers[REG_SP_ACCU]
			return &vm.registers[REG_SP_HALT], nil
		default:
			panic("unimplemented opcode")
		}
	}

	return &vm.registers[REG_SP_HALT], nil
}
