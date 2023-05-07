package introspection

type Null struct{}

func NullAPI() Null {
	return Null{}
}

// implements API
func (n Null) SendEvent(event IntrospectionEvent) {
	return
}

func (n Null) SingleStepBarrier(s IntrospectionSubject) {
	return
}
