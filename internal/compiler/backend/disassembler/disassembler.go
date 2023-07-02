package disassembler

import (
	"fmt"
	"github.com/certainty/go-braces/internal/isa"
	"io"
	"strings"
)

type Disassembler struct {
	writer io.StringWriter
}

func NewDisassembler(writer io.StringWriter) *Disassembler {
	return &Disassembler{
		writer: writer,
	}
}

func DisassModule(assemblyModule *isa.AssemblyModule) (string, error) {
	stringWriter := strings.Builder{}
	disass := NewDisassembler(&stringWriter)
	if err := disass.Disassemble(assemblyModule); err != nil {
		return "", err
	}
	return stringWriter.String(), nil
}

func (disass *Disassembler) Disassemble(assemblyModule *isa.AssemblyModule) error {
	disass.disassModuleMeta(assemblyModule.Meta)
	var address isa.InstructionAddress = 0
	var err error

	for {
		if int(address) >= len(assemblyModule.Code.Instructions) {
			break
		}
		address, err = disass.DisassInstruction(assemblyModule.Code, address)
		if err != nil {
			return err
		}
	}

	disass.writer.WriteString("\n")
	return nil
}

func (disass *Disassembler) disassModuleMeta(meta isa.AssemblyMeta) {
	moduleType := "LIB"
	if meta.Type == isa.AssemblyTypeExecutable {
		moduleType = "EXE"
	}
	disass.writer.WriteString(fmt.Sprintf("\nMOD(%s) Name: '%s' APIVersion: %x\n\n", moduleType, meta.Name, meta.ABIVersion))
}

func (disass *Disassembler) DisassInstruction(code *isa.CodeUnit, addr isa.InstructionAddress) (isa.InstructionAddress, error) {
	disass.writer.WriteString(fmt.Sprintf("0x%08x ", addr))
	instr := code.Instructions[addr]

	switch instr.Opcode {
	case isa.OP_RET, isa.OP_HALT, isa.OP_TRUE, isa.OP_FALSE:
		return disass.disassSimpleInstruction(instr, addr), nil
	case isa.OP_CONST:
		return disass.disassConstant(code, addr), nil
	default:
		return 0, fmt.Errorf("unknown opcode %d", instr.Opcode)
	}
}

func (disass *Disassembler) disassSimpleInstruction(instr isa.Instruction, addr isa.InstructionAddress) isa.InstructionAddress {
	disass.writer.WriteString(fmt.Sprintf("%-8s %s\n", disassOpCode(instr.Opcode), disassOperands(instr.Operands)))
	return addr + 1
}

func (disass *Disassembler) disassConstant(code *isa.CodeUnit, addr isa.InstructionAddress) isa.InstructionAddress {
	instr := code.Instructions[addr]
	newAddr := disass.disassSimpleInstruction(instr, addr)
	value := code.Constants[instr.Operands[0].(isa.ConstantAddress)]
	disass.writer.WriteString(fmt.Sprintf("%-8s     ^--- %v\n", "|", value))
	return newAddr
}

func disassOperands(operands []isa.Operand) string {
	ops := make([]string, len(operands))
	for _, operand := range operands {
		ops = append(ops, disassOperand(operand))
	}
	return strings.Join(ops, " ")
}

func disassOperand(operand isa.Operand) string {
	switch operand.(type) {
	case isa.Register:
		return fmt.Sprintf("$%-4d", operand)
	case isa.InstructionAddress:
		return fmt.Sprintf("@%08x ", operand)
	case isa.ConstantAddress:
		return fmt.Sprintf("%%%08x ", operand)
	case isa.Value:
		return fmt.Sprintf("%v", operand)
	default:
		return ""
	}
}

func disassOpCode(code isa.OpCode) string {
	switch code {
	case isa.OP_TRUE:
		return "TRUE"
	case isa.OP_FALSE:
		return "FALSE"
	case isa.OP_RET:
		return "RET"
	case isa.OP_HALT:
		return "HALT"
	case isa.OP_CONST:
		return "CONST"
	default:
		panic("Unknown opcode")
	}
}
