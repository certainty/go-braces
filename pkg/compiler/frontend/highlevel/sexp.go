package highlevel

var printLocation = true

type Sexp interface {
	Sexp() string
}

func SexpString(s Sexp, hideLocation bool) string {
	if hideLocation {
		printLocation = false
		defer func() { printLocation = true }()
	}

	return s.Sexp()
}

func PrintLocation() bool {
	return printLocation
}
