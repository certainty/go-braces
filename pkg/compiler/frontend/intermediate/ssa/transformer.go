package ssa

import (
	"fmt"
	"log"

	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
)

type Transformer struct {
	instrumentation compiler_introspection.Instrumentation
	Builder
}

func NewTransformer(instrumentation compiler_introspection.Instrumentation) *Transformer {
	return &Transformer{
		instrumentation: instrumentation,
		Builder:         NewBuilder(),
	}
}

func (t *Transformer) Transform(module ir.Module) (*Module, error) {
	t.instrumentation.EnterPhase(compiler_introspection.CompilationPhaseSSA)
	defer t.instrumentation.LeavePhase(compiler_introspection.CompilationPhaseSSA)

	ssaModule := &Module{Name: module.Name, Declarations: make([]Declaration, 0)}

	for _, decl := range module.Declarations {
		ssaDecl, err := t.TransformDeclaration(decl)
		if err != nil {
			return nil, err
		}
		ssaModule.Declarations = append(ssaModule.Declarations, ssaDecl)
	}
	return ssaModule, nil
}

func (t *Transformer) TransformDeclaration(decl ir.Declaration) (Declaration, error) {
	log.Printf("Transforming declaration %v", decl)
	switch decl := decl.(type) {
	case ir.ProcDecl:
		return t.TransformProc(&decl)
	default:
		return nil, fmt.Errorf("unknown declaration type %T", decl)
	}
}

func (t *Transformer) TransformProc(proc *ir.ProcDecl) (ProcDecl, error) {
	log.Printf("Transforming proc %v", proc)
	ssaProc := ProcDecl{irDecl: *proc, Blocks: make([]*BasicBlock, 0)}

	for _, block := range proc.Blocks {
		ssaBlock, err := t.TransformBlock(&block)
		if err != nil {
			return ProcDecl{}, err
		}
		ssaProc.Blocks = append(ssaProc.Blocks, ssaBlock)
	}
	return ssaProc, nil
}

func (t *Transformer) TransformBlock(block *ir.BlockExpr) (*BasicBlock, error) {
	log.Printf("Transforming block %v", block)

	blockBuilder := t.BlockBuilder(block.Label)

	for _, stmt := range block.Statements {
		_, _, err := t.TransformStatement(stmt, blockBuilder)
		if err != nil {
			return nil, err
		}

		log.Printf("Statements: %v", blockBuilder.block)
	}

	return blockBuilder.Close(), nil
}

func (t *Transformer) TransformStatement(stmt ir.Statement, block *BasicBlockBuilder) (Variable, bool, error) {
	log.Printf("Transforming statement %v", stmt)

	switch stmt := stmt.(type) {
	case ir.ExprStatement:
		return t.TransformExpr(stmt.Expr, block)
	default:
		return Variable{}, false, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (t *Transformer) TransformExpr(expr ir.Expression, block *BasicBlockBuilder) (Variable, bool, error) {
	log.Printf("Transforming expression %v", expr)
	switch expr := expr.(type) {
	case ir.AtomicLitExpr:
		variable := block.AddAssignment(t.AtomicLitExpr(expr))
		return variable, true, nil

	case ir.BinaryExpr:
		left, hasVar, err := t.TransformExpr(expr.Left, block)
		if err != nil {
			return Variable{}, false, err
		}
		if !hasVar {
			return Variable{}, false, fmt.Errorf("expected variable in left-hand side of binary expression")
		}

		right, hasVar, err := t.TransformExpr(expr.Right, block)
		if err != nil {
			return Variable{}, false, err
		}
		if !hasVar {
			return Variable{}, false, fmt.Errorf("expected variable in right-hand side of binary expression")
		}

		variable := block.AddAssignment(t.BinaryExpr(expr, left, right))
		return variable, true, nil
	default:
		return Variable{}, false, fmt.Errorf("unknown expression type: %T", expr)
	}
}
