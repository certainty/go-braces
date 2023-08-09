package intermediate

import (
	"fmt"

	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	ir "github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/ast"
	"github.com/certainty/go-braces/pkg/compiler/frontend/intermediate/types"

	hl "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/ast"
	hltypes "github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/types"
)

type Context struct {
	Types   *hltypes.TypeUniverse
	Origin  token.Origin
	Module  *ir.Module
	builder *ir.Builder
}

func Lower(origin token.Origin, theAst *hl.Source, tpeUniverse *hltypes.TypeUniverse) (*ir.Module, error) {
	ctx := &Context{
		Origin:  origin,
		Types:   tpeUniverse,
		Module:  &ir.Module{},
		builder: ir.NewBuilder(),
	}
	return ctx.lower(theAst)
}

func (ctx *Context) lower(source *hl.Source) (*ir.Module, error) {
	for _, node := range source.Declarations {
		switch node := node.(type) {
		case hl.ProcDecl:
			proc, err := ctx.lowerProcDecl(node)
			if err != nil {
				return nil, err
			}
			ctx.Module.Declarations = append(ctx.Module.Declarations, proc)
		}
	}
	return ctx.Module, nil
}

func (ctx *Context) lowerProcDecl(decl hl.ProcDecl) (*ir.ProcDecl, error) {
	var err error

	procType, err := ctx.typeOf(decl)
	if err != nil {
		return nil, err
	}

	loweredType, err := ctx.lowerType(procType)
	if err != nil {
		return nil, err
	}
	id := decl.ID()
	procName := ctx.builder.Label(decl.Name.Name, &id)
	proc := ctx.builder.ProcDecl(procName, loweredType.(types.Procedure), decl)

	entryBlock := ctx.builder.BlockBuilder(ctx.builder.Label("entry", nil), &id)
	var loweredStmt ast.Statement

	for _, stmt := range decl.Body.Statements {
		loweredStmt, err = ctx.lowerStatement(stmt)
		if err != nil {
			return nil, err
		}
		entryBlock.AddStatement(loweredStmt)
	}

	if loweredStmt != nil {
		switch stmt := loweredStmt.(type) {

		case *ast.ExprStatement:
			entryBlock.ReplaceLastStatement(entryBlock.ReturnStmt(stmt.Expr))
		}
	}

	proc.Blocks = append(proc.Blocks, entryBlock.Close())
	return proc, nil
}

func (ctx *Context) lowerStatement(stmt hl.Statement) (ast.Statement, error) {
	switch stmt := stmt.(type) {
	case hl.ExprStmt:
		expr, err := ctx.lowerExpr(stmt.Expr)
		if err != nil {
			return nil, err
		}
		return ctx.builder.ExprStatement(expr), nil
	default:
		return nil, fmt.Errorf("unexpected statement type: %T", stmt)
	}
}

func (ctx *Context) lowerExpr(expr hl.Expression) (ast.Expression, error) {
	switch expr := expr.(type) {
	case hl.BasicLitExpr:
		exprType, err := ctx.typeOf(expr)
		if err != nil {
			return nil, err
		}
		loweredType, err := ctx.lowerType(exprType)
		if err != nil {
			return nil, err
		}
		return ctx.builder.AtomicLit(loweredType, expr.Token, expr.ID()), nil
	case hl.BinaryExpr:
		exprType, err := ctx.typeOf(expr)
		if err != nil {
			return nil, err
		}
		loweredType, err := ctx.lowerType(exprType)
		if err != nil {
			return nil, err
		}
		left, err := ctx.lowerExpr(expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := ctx.lowerExpr(expr.Right)
		if err != nil {
			return nil, err
		}
		return ctx.builder.BinaryExpr(loweredType, expr.Op, left, right, expr.ID()), nil
	default:
		return nil, fmt.Errorf("unexpected expression type: %T", expr)
	}
}

func (ctx *Context) lowerType(tpe hltypes.Type) (types.Type, error) {
	switch tpe := tpe.(type) {
	case hltypes.Byte:
		return types.ByteType, nil
	case hltypes.Int:
		return types.IntType, nil
	case hltypes.Float:
		return types.FloatType, nil
	case hltypes.Bool:
		return types.BoolType, nil
	case hltypes.String:
		return types.StringType, nil
	case hltypes.Procedure:
		loweredParams := make([]types.Type, len(tpe.Params))
		loweredResults := make([]types.Type, len(tpe.Results))

		for i, param := range tpe.Params {
			loweredParam, err := ctx.lowerType(param)
			if err != nil {
				return nil, err
			}
			loweredParams[i] = loweredParam
		}
		for i, result := range tpe.Results {
			loweredResult, err := ctx.lowerType(result)
			if err != nil {
				return nil, err
			}
			loweredResults[i] = loweredResult
		}

		return types.Procedure{
			Params:  loweredParams,
			Results: loweredResults,
		}, nil
	case hltypes.Unit:
		return types.VoidType, nil
	default:
		return nil, fmt.Errorf("unexpected type: %T", tpe)
	}
}

func (ctx *Context) typeOf(n hl.Node) (hltypes.Type, error) {
	tpe, ok := ctx.Types.ExpressionTypes[n.ID()]
	if !ok {
		return nil, fmt.Errorf("no type for node %v", n)
	}
	return tpe, nil
}
