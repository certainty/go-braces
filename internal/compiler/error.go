package compiler

type CompilerError struct {
	Msg     string
	Details []error
}

func (e CompilerError) Error() string {
	return e.Msg
}
