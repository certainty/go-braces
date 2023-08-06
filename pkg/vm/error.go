package vm

type VmError struct {
	Msg     string
	Details []error
}

func (e VmError) Error() string {
	return e.Msg
}
