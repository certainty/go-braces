package introspection

type IntrospectionEvent interface{}
type IntrospectionSubject interface{}

type API interface {
	SendEvent(IntrospectionEvent)
	// waits for the resume of a single step
	SingleStepBarrier(IntrospectionSubject)
}
