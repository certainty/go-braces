package frame

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/main_section"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/statusbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/topbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

type Section int

const (
	SectionTopBar Section = iota
	SectionMain
	SectionEventLog
	SectionStatusBar
)

type model struct {
	width, height int
	theme         theme.Theme
	keyMap        GlobalKeyMap
	status        State
	err           error

	// sections
	// sectionTopBar    topbar.Model
	// sectionMain      main_section.Model
	// sectionEventLog  eventlog.Model
	// sectionStatusBar statusbar.Model
	sections []tea.Model

	// styles
	styleFrame       lipgloss.Style
	styleSectionMain lipgloss.Style

	client *compiler_introspection.Client
}

func New(client *compiler_introspection.Client) model {
	theme := theme.NewCatpuccinTheme()
	keyMap := NewGlobalKeyMap()

	//topbar
	sectionTopBar := topbar.New(theme, "(Go-Braces-Introspect 'Compiler)")

	//main
	sectionMain := main_section.New(theme)

	// sectionEventLog
	sectionEventLog := eventlog.New()

	// statusbar
	shortcuts := []*key.Binding{
		&keyMap.Continue,
		&keyMap.ToggleEventLog,
		&keyMap.Quit,
		&keyMap.Help,
	}
	sectionStatusBar := statusbar.New(theme, shortcuts)

	sections := make([]tea.Model, 4)
	sections[SectionTopBar] = sectionTopBar
	sections[SectionMain] = sectionMain
	sections[SectionEventLog] = sectionEventLog
	sections[SectionStatusBar] = sectionStatusBar

	return model{
		width:  100,
		height: 80,
		theme:  theme,
		keyMap: keyMap,
		status: Disconnected,
		client: client,

		sections:         sections,
		styleFrame:       lipgloss.NewStyle(),
		styleSectionMain: lipgloss.NewStyle(),
	}
}
