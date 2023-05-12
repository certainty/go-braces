package vm

import (
	"github.com/certainty/go-braces/internal/introspection/vm_introspection"
	"github.com/certainty/go-braces/internal/isa"
	"log"
)

const (
	REG_GP_COUNT = 256
	REG_SP_COUNT = 16
	REG_SP_HALT  = 0
	REG_SP_ACCU  = 1
)

type VmOptions struct {
	instrumentation vm_introspection.Instrumentation
}

type VM struct {
	instrumentation vm_introspection.Instrumentation
	registers       [REG_SP_COUNT + REG_GP_COUNT]isa.Value
	writer          *Writer
	// read only reference
	internedStrings *InternedStringTable
	code            *isa.CodeUnit
	pc              int
}

func DefaultOptions() VmOptions {
	return VmOptions{}
}

func NewVM(options VmOptions) *VM {
	internedStrings := NewInternedStringTable()

	vm := VM{
		pc:              0,
		internedStrings: internedStrings,
		writer:          NewWriter(internedStrings),
		code:            nil,
	}

	if options.instrumentation == nil {
		vm.instrumentation = vm_introspection.NewNullInstrumentation()
	} else {
		vm.instrumentation = options.instrumentation
	}

	return &vm
}

func (vm *VM) WriteValue(value isa.Value) string {
	return vm.writer.Write(value)
}

func (vm *VM) LoadModule(module *isa.AssemblyModule) {
	vm.code = module.Code
	vm.pc = 0
}

func (vm *VM) ExecuteModule(module *isa.AssemblyModule) (isa.Value, error) {
	vm.LoadModule(module)

	for vm.pc < len((*vm.code).Instructions) {
		instr := (*vm.code).Instructions[vm.pc]
		vm.pc++

		log.Printf("Executing instruction %s", instr)
		switch instr.Opcode {
		case isa.OP_TRUE:
			vm.registers[REG_SP_ACCU] = isa.BoolValue(true)
		case isa.OP_FALSE:
			vm.registers[REG_SP_ACCU] = isa.BoolValue(false)
		case isa.OP_HALT:
			vm.registers[REG_SP_HALT] = vm.registers[REG_SP_ACCU]
			return vm.registers[REG_SP_HALT], nil
		default:
			panic("unimplemented opcode")
		}
	}

	return vm.registers[REG_SP_HALT], nil
}
