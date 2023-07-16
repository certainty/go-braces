package frame

import (
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/eventlog"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/main_section"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/statusbar"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/section/topbar"
	"github.com/certainty/go-braces/cmd/braces-introspect/internal/compiler_introspector/ui/theme"
	"github.com/certainty/go-braces/pkg/introspection/compiler_introspection"
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

	//topbar
	sectionTopBar := topbar.New(theme, "(Go-Braces-Introspect 'Compiler)")

	//main
	sectionMain := main_section.New(theme, client)

	// sectionEventLog
	sectionEventLog := eventlog.New()

	globalKeyMap := common.NewKeyMap("global", AllBindings, Shortcuts)
	contextKeyMap := common.NewKeyMap("", []key.Binding{}, []key.Binding{})
	sectionStatusBar := statusbar.New(theme, globalKeyMap, contextKeyMap)

	sections := make([]tea.Model, 4)
	sections[SectionTopBar] = sectionTopBar
	sections[SectionMain] = sectionMain
	sections[SectionEventLog] = sectionEventLog
	sections[SectionStatusBar] = sectionStatusBar

	return model{
		width:  100,
		height: 80,
		theme:  theme,
		status: Disconnected,
		client: client,

		sections:         sections,
		styleFrame:       lipgloss.NewStyle(),
		styleSectionMain: lipgloss.NewStyle(),
	}
}
