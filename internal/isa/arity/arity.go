package arity

type Arity interface {
	IsSatisfiedBy(count uint) bool
}
type atLeast struct{ count uint }
type exactly struct{ count uint }

func (a atLeast) IsSatisfiedBy(count uint) bool {
	return count >= a.count
}

func (a exactly) IsSatisfiedBy(count uint) bool {
	return count == a.count
}

func AtLeast(count uint) Arity {
	return atLeast{count}
}

func Exactly(count uint) Arity {
	return exactly{count}
}
