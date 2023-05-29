package frame

import "github.com/charmbracelet/bubbles/key"

var (
	KeyQuit           = key.NewBinding(key.WithKeys("Q"), key.WithHelp("[Q]uit", "Quit the program"))
	KeyHelp           = key.NewBinding(key.WithKeys("H"), key.WithHelp("[H]elp", "Open help dialog"))
	KeyToggleEventLog = key.NewBinding(key.WithKeys("E"), key.WithHelp("[E]vents", "Toggle event log"))

	AllBindings = []key.Binding{KeyQuit, KeyHelp, KeyToggleEventLog}
	Shortcuts   = []key.Binding{KeyToggleEventLog, KeyHelp, KeyQuit}
)
