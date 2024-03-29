package disassembler

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/certainty/go-braces/pkg/shared/isa"
)

type Disassembler struct {
	writer io.StringWriter
	indent string
}

func NewDisassembler(writer io.StringWriter) *Disassembler {
	return &Disassembler{
		writer: writer,
		indent: "",
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
	for _, function := range assemblyModule.Functions {
		disass.dissassFunction(&function)
	}
	for _, closure := range assemblyModule.Closures {
		disass.disassClosure(&closure)
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

func (disass *Disassembler) dissassFunction(function *isa.Function) {
	disass.writer.WriteString(fmt.Sprintf("@%s:\n", function.Label))
	disass.indent = "  "
	disass.disassCodeUnit(&function.Code)
	disass.indent = ""
	disass.writer.WriteString("\n")
}

func (disass *Disassembler) disassClosure(closure *isa.Closure) {
	disass.writer.WriteString(fmt.Sprintf("@%s:\n", closure.Function.Label))
	disass.indent = "  "
	disass.disassCodeUnit(&closure.Function.Code)
	disass.indent = ""
	disass.writer.WriteString("\n")
}

func (disass *Disassembler) disassCodeUnit(code *isa.CodeUnit) {
	var err error
	instructionCount := len(code.Instructions)
	for addr := isa.InstructionAddress(0); int(addr) < instructionCount; {
		addr, err = disass.DisassInstruction(code, addr)
		if err != nil {
			log.Fatalf("disassembler: %v", err)
		}
	}
}

func (disass *Disassembler) DisassInstruction(code *isa.CodeUnit, addr isa.InstructionAddress) (isa.InstructionAddress, error) {
	disass.writer.WriteString(fmt.Sprintf("%s0x%08x ", disass.indent, addr))
	instr := code.Instructions[addr]

	switch instr.Opcode {
	case isa.OP_RET, isa.OP_STORE, isa.OP_HALT, isa.OP_ADD, isa.OP_ADDI, isa.OP_SUB, isa.OP_SUBI, isa.OP_MUL, isa.OP_DIV:
		return disass.disassSimpleInstruction(instr, addr), nil
	case isa.OP_LOAD:
		return disass.disassConstant(code, addr), nil
	default:
		return 0, fmt.Errorf("unknown opcode %d", instr.Opcode)
	}
}

func (disass *Disassembler) disassSimpleInstruction(instr isa.Instruction, addr isa.InstructionAddress) isa.InstructionAddress {
	disass.writer.WriteString(fmt.Sprintf("%s%-8s %s\n", disass.indent, disassOpCode(instr.Opcode), disassOperands(instr.Opcode, instr.Operands)))
	return addr + 1
}

func (disass *Disassembler) disassConstant(code *isa.CodeUnit, addr isa.InstructionAddress) isa.InstructionAddress {
	instr := code.Instructions[addr]
	newAddr := disass.disassSimpleInstruction(instr, addr)
	value := code.Constants[instr.Operands[1]]
	disass.writer.WriteString(fmt.Sprintf("  %-8s     %-17s^--- %v\n", "|", "", value))
	return newAddr
}

func disassOperands(op isa.OpCode, operands []isa.Operand) string {
	switch op {
	case isa.OP_LOAD:
		return fmt.Sprintf("%s, %s", asRegister(operands[0]), asConstAddress(operands[1]))
	case isa.OP_STORE:
		return fmt.Sprintf("%s, %s", asRegister(operands[0]), asRegister(operands[1]))
	case isa.OP_RET:
		return asRegister(operands[0])
	case isa.OP_ADD, isa.OP_SUB, isa.OP_MUL, isa.OP_DIV:
		return fmt.Sprintf("%s, %s, %s", asRegister(operands[0]), asRegister(operands[1]), asRegister(operands[2]))
	case isa.OP_ADDI, isa.OP_SUBI:
		return fmt.Sprintf("%s, %s, %s", asRegister(operands[0]), asRegister(operands[1]), asImmediate(operands[2]))
	default:
		return ""
	}
}

func asRegister(operand isa.Operand) string {
	return fmt.Sprintf("$%-4d", operand)
}

func asConstAddress(operand isa.Operand) string {
	return fmt.Sprintf("%%%08x ", operand)
}

func asImmediate(operand isa.Operand) string {
	return fmt.Sprintf("%d", operand)
}

var opCodeMap = map[isa.OpCode]string{
	isa.OP_LOAD:  "LOAD",
	isa.OP_LOADI: "LOADI",
	isa.OP_STORE: "STORE",
	isa.OP_ADD:   "ADD",
	isa.OP_ADDI:  "ADDI",
	isa.OP_SUB:   "SUB",
	isa.OP_SUBI:  "SUBI",
	isa.OP_MUL:   "MUL",
	isa.OP_DIV:   "DIV",
	isa.OP_MOD:   "MOD",
	isa.OP_AND:   "AND",
	isa.OP_OR:    "OR",
	isa.OP_RET:   "RET",
	isa.OP_HALT:  "HALT",
}

func disassOpCode(code isa.OpCode) string {
	if op, ok := opCodeMap[code]; ok {
		return op
	}
	panic("Unknown opcode")
}
