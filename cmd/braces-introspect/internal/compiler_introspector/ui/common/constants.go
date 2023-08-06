package common

type Mode int

const (
	NormalMode Mode = iota
	WaitingMode
	CompileMode
	SingleStepMode
)

func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "Normal"
	case WaitingMode:
		return "Waiting"
	case CompileMode:
		return "Compile"
	case SingleStepMode:
		return "SingleStep"
	default:
		return "Unknown"
	}
}

type RequestStatus int

const (
	NoRequest RequestStatus = iota
	RequestSent
	RequestDone
)
