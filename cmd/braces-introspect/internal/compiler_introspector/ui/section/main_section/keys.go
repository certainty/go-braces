package main_section

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/charmbracelet/bubbles/key"
)

var (
	KeyContinue = key.NewBinding(key.WithKeys("C"), key.WithHelp("[C]ontinue", "Continue"))

	AllBindings   = []key.Binding{KeyContinue}
	Shortcuts     = []key.Binding{KeyContinue}
	CompileKeyMap = common.NewKeyMap("compile", AllBindings, Shortcuts)
)
