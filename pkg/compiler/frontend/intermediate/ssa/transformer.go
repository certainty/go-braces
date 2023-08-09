package ssa

import (
	"fmt"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
	log "github.com/sirupsen/logrus"
)

type Transformer struct {
	instrumentation compiler_introspection.Instrumentation
	*ir.Builder
}

func NewTransformer(instrumentation compiler_introspection.Instrumentation) *Transformer {
	return &Transformer{
		instrumentation: instrumentation,
		Builder:         ir.NewBuilder(),
	}
}

func (t *Transformer) Transform(module *ir.Module) error {
	ir.Walk(t, module)
	return nil
}

func (t *Transformer) Leave(node ir.Node) {}

func (t *Transformer) Enter(node ir.Node) bool {
	switch n := node.(type) {
	case *ir.ProcDecl:
		err := t.TransformProc(n)
		if err != nil {
			log.Warnf("Error transforming proc %v: %v", n, err)
		}
		return false
	default:
		// nothing
	}
	return true
}

func (t *Transformer) TransformProc(proc *ir.ProcDecl) error {
	log.Debugf("Transforming proc %v", proc)
	ssaBlocks := make([]*ir.BasicBlock, 0)

	for _, block := range proc.Blocks {
		ssaBlock, err := t.TransformBlock(block)
		if err != nil {
			return err
		}
		ssaBlocks = append(ssaBlocks, ssaBlock)
	}
	proc.SSABlocks = ssaBlocks
	return nil
}

func (t *Transformer) TransformBlock(block *ir.BasicBlock) (*ir.BasicBlock, error) {
	log.Debugf("Transforming block %v", block)

	blockBuilder := t.BlockBuilder(block.Label, nil)

	for _, stmt := range block.Statements {
		_, err := t.TransformStatement(stmt, blockBuilder)
		if err != nil {
			return nil, err
		}
	}

	return blockBuilder.Close(), nil
}

func (t *Transformer) TransformStatement(stmt ir.Statement, block *ir.BlockBuilder) (ir.Expression, error) {
	switch stmt := stmt.(type) {
	case *ir.ExprStatement:
		return t.TransformExpr(stmt.Expr, block)
	case *ir.ReturnStmt:
		expr, err := t.TransformExpr(stmt.Value, block)
		if err != nil {
			return nil, err
		}
		block.AddStatement(t.ReturnStmt(expr))
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (t *Transformer) TransformExpr(expr ir.Expression, block *ir.BlockBuilder) (ir.Expression, error) {
	switch expr := expr.(type) {
	case *ir.AtomicLitExpr:
		return expr, nil

	case *ir.BinaryExpr:
		leftExpr, err := t.TransformExpr(expr.Left, block)
		if err != nil {
			return nil, err
		}

		rightExpr, err := t.TransformExpr(expr.Right, block)
		if err != nil {
			return nil, err
		}

		variable := t.Variable("temp")
		block.AddAssignment(variable, t.BinaryExpr(expr.Type, expr.Op, leftExpr, rightExpr, expr.ID()))
		return variable, nil
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}
