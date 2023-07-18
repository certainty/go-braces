package ast

type Visitor interface {
	Enter(node Node) bool // return true to continue visiting
	Leave(node Node)
}

func walkDeclarations(v Visitor, decls []Declaration) {
	for _, decl := range decls {
		Walk(v, decl)
	}
}

func walkStatements(v Visitor, stmts []Statement) {
	for _, stmt := range stmts {
		Walk(v, stmt)
	}
}

func Walk(v Visitor, node Node) {
	cont := v.Enter(node)
	defer v.Leave(node)

	if !cont {
		return
	}

	switch n := node.(type) {
	case *Source:
		walkDeclarations(v, n.Declarations)
	case *BadDecl, *TypeSpec, *Identifier:
		/// nothing to do
	case *ProcDecl:
		Walk(v, n.Name)
		Walk(v, n.Body)
	case *Field:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *BadStmt:
		// nothing
	case *ExprStmt:
		Walk(v, n.Expr)
	case *BadExpr, *BasicLitExpr:
	//nothing
	case *BlockExpr:
		walkStatements(v, n.Statements)
	case *ParenExpr:
		Walk(v, n.Expr)
	case *UnaryExpr:
		Walk(v, n.Expr)
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	default:
		// nothing
	}
}
