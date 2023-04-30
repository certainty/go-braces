package reader

import "fmt"

type ReadError struct {
	Msg string
	pos Position
}

func (e ReadError) Error() string {
	return fmt.Sprintf("%s at %d:%d", e.Msg, e.pos.Line, e.pos.Col)
}
