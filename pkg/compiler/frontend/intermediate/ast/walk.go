package ast

type Visitor interface {
	Enter(node Node) bool // return true to continue visiting
	Leave(node Node)
}

func Walk(v Visitor, node Node, visitSSABlocks bool) {
	cont := v.Enter(node)
	defer v.Leave(node)
	if !cont {
		return
	}

	switch n := node.(type) {
	case *Module:
		for _, decl := range n.Declarations {
			Walk(v, decl, visitSSABlocks)
		}
	case *BasicBlock:
		for _, stmt := range n.Statements {
			Walk(v, stmt, visitSSABlocks)
		}
	case *BinaryExpr:
		Walk(v, n.Left, visitSSABlocks)
		Walk(v, n.Right, visitSSABlocks)

	case *AtomicLitExpr:
		// nothing
	case *ExprStatement:
		Walk(v, n.Expr, visitSSABlocks)
	case *Phi:
	// nothing
	case *Variable:
		// nothing
	case *AssignStmt:
		Walk(v, n.Variable, visitSSABlocks)
		Walk(v, n.Expr, visitSSABlocks)
	case *ReturnStmt:
		Walk(v, n.Value, visitSSABlocks)
	case *ProcDecl:
		Walk(v, n.Name, visitSSABlocks)

		if visitSSABlocks {
			for _, block := range n.SSABlocks {
				Walk(v, block, visitSSABlocks)
			}
		} else {
			for _, block := range n.Blocks {
				Walk(v, block, visitSSABlocks)
			}
		}
	default: //nothing
	}
}
