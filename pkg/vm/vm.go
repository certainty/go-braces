package vm

import (
	"fmt"
	"github.com/certainty/go-braces/pkg/compiler/backend/disassembler"
	"github.com/certainty/go-braces/pkg/introspection/vm_introspection"
	"github.com/certainty/go-braces/pkg/shared/isa"
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

		switch instr.Opcode {
		case isa.OP_RET:
			return vm.registers[instr.Operands[0]], nil
		case isa.OP_MUL:
			target := instr.Operands[0]
			left := instr.Operands[1]
			right := instr.Operands[2]
			vm.registers[target] = isa.Int(vm.registers[left].(int) * vm.registers[right].(int))
		case isa.OP_ADD:
			target := instr.Operands[0]
			left := instr.Operands[1]
			right := instr.Operands[2]
			// TODO: make sure correct values are used
			vm.registers[target] = isa.Int(vm.registers[left].(int) + vm.registers[right].(int))
		case isa.OP_ADDI:
			target := instr.Operands[0]
			left := instr.Operands[1]
			right := int(instr.Operands[2])
			vm.registers[target] = (vm.registers[left].(int) + right)
		case isa.OP_LOAD:
			register := instr.Operands[0]
			address := instr.Operands[1]
			value, err := vm.code.ReadConstant(isa.ConstantAddress(address))
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
