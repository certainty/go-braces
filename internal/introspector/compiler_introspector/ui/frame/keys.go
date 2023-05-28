package frame

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Quit           key.Binding
	Help           key.Binding
	Continue       key.Binding
	Documentation  key.Binding
	ToggleEventLog key.Binding
}

func NewGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit:           key.NewBinding(key.WithKeys("Q"), key.WithHelp("[Q]uit", "Quit the program")),
		Help:           key.NewBinding(key.WithKeys("H"), key.WithHelp("[H]elp", "Open help dialog")),
		Documentation:  key.NewBinding(key.WithKeys("D"), key.WithHelp("[D]ocumentation", "Open documentation dialog")),
		Continue:       key.NewBinding(key.WithKeys("C"), key.WithHelp("[C]ontinue", "Continue")),
		ToggleEventLog: key.NewBinding(key.WithKeys("E"), key.WithHelp("[E]vents", "Toggle event log")),
	}
}
