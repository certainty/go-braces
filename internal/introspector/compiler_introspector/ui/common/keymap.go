package common

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Context   string
	Bindings  []key.Binding
	ShortCuts []key.Binding
}

func NewKeyMap(context string, bindings []key.Binding, shortcuts []key.Binding) KeyMap {
	return KeyMap{
		context, bindings, shortcuts,
	}
}
