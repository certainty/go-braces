package ssa

import (
	"fmt"

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
		ssaDecl, err := t.TransformDeclaration(&decl)
		if err != nil {
			return nil, err
		}
		ssaModule.Declarations = append(ssaModule.Declarations, ssaDecl)
	}
	return ssaModule, nil
}

func (t *Transformer) TransformDeclaration(decl *ir.Declaration) (Declaration, error) {
	switch decl := (*decl).(type) {
	case ir.ProcDecl:
		return t.TransformProc(&decl)
	default:
		return nil, fmt.Errorf("unknown declaration type %T", decl)
	}
}

func (t *Transformer) TransformProc(proc *ir.ProcDecl) (ProcDecl, error) {
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
	blockBuilder := t.BlockBuilder(block.Label)

	for _, stmt := range block.Statements {
		_, err := t.TransformStatement(&stmt, &blockBuilder)
		if err != nil {
			return nil, err
		}
	}

	return blockBuilder.Close(), nil
}

func (t *Transformer) TransformStatement(stmt *ir.Statement, block *BasicBlockBuilder) (*Variable, error) {
	switch stmt := (*stmt).(type) {
	case ir.ExprStatement:
		return t.TransformExpr(stmt.Expr, block)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (t *Transformer) TransformExpr(expr ir.Expression, block *BasicBlockBuilder) (*Variable, error) {
	switch expr := expr.(type) {
	case ir.AtomicLitExpr:
		variable := block.AddAssingment(t.AtomicLitExpr(expr))
		return &variable, nil
	case ir.BinaryExpr:
		left, err := t.TransformExpr(expr.Left, block)
		if err != nil {
			return nil, err
		}

		right, err := t.TransformExpr(expr.Right, block)
		if err != nil {
			return nil, err
		}

		variable := block.AddAssingment(t.BinaryExpr(expr, *left, *right))
		return &variable, nil
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}
