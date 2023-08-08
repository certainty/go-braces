package astutils

type Visitor[Node any] interface {
	Enter(node Node) bool // return true to continue visiting
	Leave(node Node)
}

type Walker[Node any] interface {
	Walk(visitor Visitor[Node], node Node)
}
