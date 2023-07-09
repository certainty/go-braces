package vm

import (
	"fmt"
	"log"

	"github.com/certainty/go-braces/internal/compiler/backend/disassembler"
	"github.com/certainty/go-braces/internal/introspection/vm_introspection"
	"github.com/certainty/go-braces/internal/isa"
)

type VmOptions struct {
	instrumentation vm_introspection.Instrumentation
}

type VM struct {
	instrumentation vm_introspection.Instrumentation
	registers       [isa.REG_SP_COUNT + isa.REG_GP_COUNT]isa.Value
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

func (vm *VM) LoadModule(module *isa.AssemblyModule) error {
	if module.EntryPoint < 0 {
		return fmt.Errorf("invalid entry point")
	}

	vm.code = &module.Functions[module.EntryPoint].Code
	vm.pc = 0

	return nil
}

func (vm *VM) ExecuteModule(module *isa.AssemblyModule) (isa.Value, error) {
	if err := vm.LoadModule(module); err != nil {
		return nil, err
	}

	fmt.Print(disassembler.DisassModule(module))

	for vm.pc < len((*vm.code).Instructions) {
		instr := (*vm.code).Instructions[vm.pc]
		vm.pc++

		log.Printf("Executing instruction %s", instr)
		switch instr.Opcode {
		case isa.OP_TRUE:
			register := instr.Operands[0].(isa.Register)
			vm.registers[register] = isa.UInt8(1)
		case isa.OP_FALSE:
			register := instr.Operands[0].(isa.Register)
			vm.registers[register] = isa.UInt8(0)
		case isa.OP_CONST:
			address := instr.Operands[0].(isa.ConstantAddress)
			register := instr.Operands[1].(isa.Register)
			value, err := vm.code.ReadConstant(address)
			if err != nil {
				vm.panic("invalid constant")
			}
			vm.registers[register] = value
		case isa.OP_HALT:
			vm.registers[isa.REG_SP_HALT] = vm.registers[isa.REG_SP_ACCU]
			return vm.registers[isa.REG_SP_HALT], nil
		default:
			panic("unimplemented opcode")
		}
	}

	return vm.registers[isa.REG_SP_HALT], nil
}

func (vm *VM) panic(message string) {
	// just panic for now
	panic(message)
}
