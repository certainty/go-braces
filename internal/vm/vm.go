package vm

import (
	"log"

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

	log.Printf("Executing module %v", module.Code.Instructions)

	for vm.pc < len((*vm.code).Instructions) {
		instr := (*vm.code).Instructions[vm.pc]
		vm.pc++

		log.Printf("PC:%d OP_CODE: %d", vm.pc, instr.Opcode)
		log.Printf("Registers: %v", vm.registers)

		switch instr.Opcode {
		case isa.OP_TRUE:
			log.Printf("OP_TRUE")
			vm.registers[REG_SP_ACCU] = isa.BoolValue(true)
		case isa.OP_HALT:
			vm.registers[REG_SP_HALT] = vm.registers[REG_SP_ACCU]
			break
		default:
			panic("unimplemented opcode")
		}
	}

	log.Printf("Halted with %v", vm.registers[REG_SP_HALT])

	return &vm.registers[REG_SP_HALT], nil
}
