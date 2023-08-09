package main

import (
	"fmt"
	"os"

	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"

	log "github.com/sirupsen/logrus"
)

type Parent interface {
	parent()
}
type Child interface {
	Parent
	child()
}

type Block struct {
	Child
	Children []Child
}

type ConcreteChild struct {
	Child
	Value int
}

func (c *ConcreteChild) parent() {}
func (c *ConcreteChild) child()  {}
func (b *Block) parent()         {}
func (b *Block) child()          {}

type Visitor interface {
	Enter(node Parent) bool
	Leave(node Parent)
}

func Walk(v Visitor, node Parent) {
	cont := v.Enter(node)
	defer v.Leave(node)
	if !cont {
		return
	}

	switch n := node.(type) {
	case *Block:
		for _, child := range n.Children {
			Walk(v, child)
		}
	case *ConcreteChild:
		// nothing
	default:
		log.Debugf("Skipping SSA transformation for %T", node)
	}
}

type AddingVisitor struct{}

func (v *AddingVisitor) Enter(node Parent) bool {
	switch n := node.(type) {
	case *ConcreteChild:
		n.Value += 100
	default:
		// nothing
	}
	return true
}

func (v *AddingVisitor) Leave(node Parent) {}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	compiler_introspection.RegisterTypes()

	tst := &Block{
		Children: make([]Child, 0),
	}
	tst.Children = append(tst.Children, &ConcreteChild{Value: 1})
	tst.Children = append(tst.Children, &ConcreteChild{Value: 2})
	tst.Children = append(tst.Children, &ConcreteChild{Value: 3})

	v := &AddingVisitor{}

	Walk(v, tst)

	for _, child := range tst.Children {
		log.Infof("child: %d", child.(*ConcreteChild).Value)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
