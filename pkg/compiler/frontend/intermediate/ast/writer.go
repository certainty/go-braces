package ir

import (
	"fmt"
	"log"
	"strings"

	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"
)

type IRWriter struct {
	module *Module
}

func NewIRWriter(module *Module) *IRWriter {
	return &IRWriter{
		module: module,
	}
}

func (w *IRWriter) WriteProcedure(proc *Procedure) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("define %s @%s(", w.writeType(proc.tpe), proc.Name))
	args := make([]string, len(proc.Args))
	for i, a := range proc.Args {
		args[i] = fmt.Sprintf("%s %s", w.writeType(a.tpe), w.writeRegister(a.Register))
	}
	result.WriteString(strings.Join(args, ", "))
	result.WriteString(") {\n")
	for _, b := range proc.Blocks {
		result.WriteString(w.WriteBlock(b))
	}
	result.WriteString("}\n")
	return result.String()
}

func (w *IRWriter) WriteBlock(b *BasicBlock) string {
	result := fmt.Sprintf("  %s\n", w.writeLabel(b.Label))
	for _, i := range b.Instructions {
		result += fmt.Sprintf("    %s\n", w.WriteInstruction(i))
	}
	return result
}

func (w *IRWriter) WriteInstruction(i Instruction) string {
	switch i := i.(type) {
	case SimpleInstruction:
		ops := make([]string, len(i.Operands))
		for i, o := range i.Operands {
			ops[i] = w.writeOperand(o)
		}
		return fmt.Sprintf("%s = %s %s %s", w.writeRegister(i.Register), w.writeOperation(i.Operation), w.writeType(i.tpe), strings.Join(ops, ", "))
	case AssignmentInstruction:
		return fmt.Sprintf("%s = %s %s", w.writeRegister(i.Register), w.writeType(i.tpe), w.writeOperand(i.Operand))
	case ReturnInstruction:
		return fmt.Sprintf("ret %s %s", w.writeType(i.tpe), w.writeRegister(i.Register))
	default:
		log.Printf("Unknown instruction: %s", i)
	}
	return ""
}

func (w *IRWriter) writeLabel(l Label) string {
	return fmt.Sprintf("%s: ", l)
}

func (w *IRWriter) writeRegister(r Register) string {
	return fmt.Sprintf("%%%d", r)
}

func (w *IRWriter) writeType(t types.Type) string {
	switch t.(type) {
	case types.Bool:
		return "bool"
	case types.Int:
		return "int"
	case types.UInt:
		return "uint"
	case types.Float:
		return "float"
	case types.String:
		return "string"
	default:
		panic("unknown type")
	}
}

func (w *IRWriter) writeOperation(op Operation) string {
	switch op {
	case Add:
		return "add"
	case Sub:
		return "sub"
	case Mul:
		return "mul"
	case Div:
		return "div"
	case Or:
		return "or"
	case And:
		return "and"
	case Xor:
		return "xor"
	default:
		panic("unknown operation")
	}
}

func (w *IRWriter) writeOperand(o Operand) string {
	switch o := o.(type) {
	case Label:
		return string(o)
	case Register:
		return fmt.Sprintf("%%%d", o)
	case Literal:
		return fmt.Sprintf("%v", o)
	default:
		panic("unknown operand type")
	}
}

func (w *IRWriter) Write() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("__module: %s\n", w.module.Name))

	for _, proc := range w.module.Procedures {
		sb.WriteString(w.WriteProcedure(&proc))
	}

	return sb.String()
}
