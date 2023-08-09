package ast

type Visitor interface {
	Enter(node Node) bool // return true to continue visiting
	Leave(node Node)
}

func Walk(v Visitor, node Node) {
	cont := v.Enter(node)
	defer v.Leave(node)
	if !cont {
		return
	}

	switch n := node.(type) {
	case *Module:
		for _, decl := range n.Declarations {
			Walk(v, decl)
		}
	case *BasicBlock:
		for _, stmt := range n.Statements {
			Walk(v, stmt)
		}
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *AtomicLitExpr:
		// nothing
	case *ExprStatement:
		Walk(v, n.Expr)
	case *Phi:
	// nothing
	case *Variable:
		// nothing
	case *AssignStmt:
		Walk(v, n.Expr)
	case *ReturnStmt:
		Walk(v, n.Value)
	case *ProcDecl:
		Walk(v, n.Name)
		for _, block := range n.Blocks {
			Walk(v, block)
		}
	default: //nothing
	}
}
