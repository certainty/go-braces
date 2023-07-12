package ir

import (
	"fmt"
	"log"
	"strings"
)

type IRWriter struct {
	module *Module
}

func NewIRWriter(module *Module) *IRWriter {
	return &IRWriter{
		module: module,
	}
}

func (w *IRWriter) WriteFunction(f *Function) string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("define %s @%s(", w.writeType(f.tpe), f.Name))
	args := make([]string, len(f.Args))
	for i, a := range f.Args {
		args[i] = fmt.Sprintf("%s %s", w.writeType(a.tpe), w.writeRegister(a.Register))
	}
	result.WriteString(strings.Join(args, ", "))
	result.WriteString(") {\n")
	for _, b := range f.Blocks {
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

func (w *IRWriter) writeType(t Type) string {
	switch t.(type) {
	case Bool:
		return "bool"
	case Int:
		return "i64"
	case UInt:
		return "u64"
	case Float:
		return "f32"
	case String:
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

	for _, f := range w.module.Functions {
		sb.WriteString(w.WriteFunction(f))
	}

	return sb.String()
}
