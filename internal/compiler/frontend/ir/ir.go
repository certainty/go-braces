package ir

import "github.com/certainty/go-braces/internal/compiler/frontend/parser"

type IRInstruction interface{}

type IRConstant struct{}
type IRLabel struct{}
type IRGlobalRef struct{}
type IRSet struct{}
type IRClosure struct{}
type IRCall struct{}
type IRTailCall struct{}
type IRBranch struct{}
type IRRet struct{}

type IR struct{}

func LowerToIR(coreAst *parser.CoreAST) (*IR, error) {
	return nil, nil
}
