package components

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"
)

type StatusBarModel struct {
	ContainerWidth  int
	ContainerHeight int
	theme           theme.Theme
	impl            statusbar.Bubble
	Phase           string
	Shortcuts       []key.Binding
	RequestState    IntrospectionRequestState
	RequestSpinner  spinner.Model
	Errors          string
}

func InitialStatusBarModel(theme theme.Theme, shortcuts []key.Binding) StatusBarModel {
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: theme.Colors.Background,
			Background: theme.Colors.Blue,
		},
		statusbar.ColorConfig{
			Foreground: theme.Colors.Green,
			Background: theme.Colors.Background,
		},
		statusbar.ColorConfig{
			Foreground: theme.Colors.Text,
			Background: theme.Colors.Background,
		},
		statusbar.ColorConfig{
			Foreground: theme.Colors.Text,
			Background: theme.Colors.Background,
		},
	)

	requestSpinner := spinner.New()
	requestSpinner.Spinner = spinner.Dot

	return StatusBarModel{
		theme:          theme,
		Phase:          "waiting",
		Errors:         "",
		Shortcuts:      shortcuts,
		RequestSpinner: requestSpinner,
		impl:           sb,
	}
}

func (m StatusBarModel) RenderShortCut(shortcut key.Binding) string {
	return m.theme.Statusbar.Copy().Faint(!shortcut.Enabled()).PaddingRight(2).Render(shortcut.Help().Key)
}

func (m StatusBarModel) Init() tea.Cmd {
	return m.RequestSpinner.Tick
}

func (m StatusBarModel) Update(msg tea.Msg) (StatusBarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.RequestSpinner, cmd = m.RequestSpinner.Update(msg)
		return m, cmd
	case messages.IntrospectionEventMsg:
		switch evt := msg.Event.(type) {
		case compiler_introspection.EventEnterPhase:
			m.Phase = string(evt.Phase)
		case compiler_introspection.EventBeginCompileModule:
			m.Phase = "compile"
		case compiler_introspection.EventEndCompileModule:
			m.Phase = "waiting"
		}
	}
	return m, nil
}

func (m StatusBarModel) View() string {
	renderedShortcuts := []string{
		"\uF11C  ",
	}

	for _, shortcut := range m.Shortcuts {
		renderedShortcuts = append(renderedShortcuts, m.RenderShortCut(shortcut))
	}

	var errors string
	if m.Errors != "" {
		errors = m.theme.Statusbar.Copy().Foreground(m.theme.Colors.Error).Render(m.Errors)
	} else {
		errors = m.theme.Statusbar.Copy().Foreground(m.theme.Colors.Success).Render("no errors")
	}

	connectionStatus := ""
	if m.RequestState != NoRequest {
		connectionStatus = m.RequestSpinner.View()
	}

	m.impl.SetSize(m.ContainerWidth)
	m.impl.SetContent(m.Phase, errors, lipgloss.JoinHorizontal(lipgloss.Top, renderedShortcuts...), connectionStatus)

	return m.impl.View()
}
