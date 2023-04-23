package value

// Use a factory to create scheme values
type Factory struct {
	true        Value
	false       Value
	nil         Value
	unspecified Value
}

func NewFactory() *Factory {
	return &Factory{
		true:        true,
		false:       false,
		nil:         Nil,
		unspecified: Unspecified,
	}
}

func (f *Factory) Bool(v bool) *Value {
	if v {
		return &f.true
	} else {
		return &f.false
	}
}

func (f *Factory) Nil() *Value {
	return &f.nil
}

func (f *Factory) Unspecified() *Value {
	return &f.unspecified
}
