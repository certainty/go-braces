package vm_introspection

type Instrumentation interface{}

type NullInstrumentation struct{}

func NewNullInstrumentation() *NullInstrumentation {
	return &NullInstrumentation{}
}
