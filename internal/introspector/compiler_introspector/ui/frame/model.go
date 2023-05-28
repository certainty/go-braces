package frame

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/statusbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/topbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	Initial State = iota
	Disconnected
	Connecting
	Connected
	Compiling
	Error
)

type model struct {
	width, height int
	theme         theme.Theme
	keyMap        GlobalKeyMap
	status        State
	err           error

	// sections
	topBar      topbar.Model
	mainSection main_section.Model
	eventLog    eventlog.Model
	statusBar   statusbar.Model

	// styles
	frameStyle         lipgloss.Style
	mainContainerStyle lipgloss.Style

	client *compiler_introspection.Client
}

func New(client *compiler_introspection.Client) model {
	theme := theme.NewCatpuccinTheme()
	keyMap := NewGlobalKeyMap()

	//topbar

	//main

	// statusbar
	shortcuts := []*key.Binding{
		&keyMap.ToggleEventLog,
		&keyMap.Quit,
		&keyMap.Help,
	}
	topBar := topbar.New(theme, "(Go-Braces-Introspect 'Compiler)")
	mainSection := main_section.New(theme)
	eventlog := eventlog.New()
	statusBar := statusbar.New(theme, shortcuts)

	return model{
		width:  100,
		height: 80,
		theme:  theme,
		keyMap: keyMap,
		status: Disconnected,
		client: client,

		topBar:      topBar,
		mainSection: mainSection,
		eventLog:    eventlog,
		statusBar:   statusBar,

		frameStyle:         lipgloss.NewStyle(),
		mainContainerStyle: lipgloss.NewStyle(),
	}
}
