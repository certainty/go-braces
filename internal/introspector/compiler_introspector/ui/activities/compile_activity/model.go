package compile_activity

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/certainty/go-braces/internal/compiler/location"
// 	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/activities"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components/phase_indicator"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
// 	"github.com/charmbracelet/bubbles/viewport"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/philistino/teacup/code"
// )

// type ReadModel struct {
// 	containerHeight int
// 	containerWidth  int
// 	code            viewport.Model
// 	sourceCode      string
// 	theme           theme.Theme
// }

// func NewReadModel(theme theme.Theme) ReadModel {
// 	return ReadModel{
// 		code:  viewport.New(0, 0),
// 		theme: theme,
// 	}
// }

// func (m *ReadModel) setSourceCode(origin location.Origin, sourceCode string) error {
// 	highlightedCode, err := code.Highlight(sourceCode, "rkt", "catppuccin-macchiato")
// 	// now prepend line numbers
// 	m.sourceCode = m.prependLineNumbers(highlightedCode)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (m ReadModel) prependLineNumbers(code string) string {
// 	lineNumberStyle := lipgloss.NewStyle().Foreground(m.theme.Colors.Overlay0)
// 	lines := strings.Split(code, "\n")
// 	var numberedLines []string
// 	for i, line := range lines {
// 		number := lineNumberStyle.Render(fmt.Sprintf("%d", i+1))
// 		numberedLines = append(numberedLines, fmt.Sprintf("%s %s", number, line))
// 	}
// 	return strings.Join(numberedLines, "\n")
// }

// func (m ReadModel) Init() tea.Cmd {
// 	return nil
// }

// func (m ReadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd

// 	m.code, cmd = m.code.Update(msg)
// 	return m, cmd
// }

// func (m ReadModel) View() string {
// 	border := lipgloss.NormalBorder()
// 	content := lipgloss.NewStyle().Width(m.code.Width - 2).Height(m.code.Height - 2).PaddingLeft(1).Border(border).Render(m.sourceCode)
// 	m.code.SetContent(content)
// 	return m.code.View()
// }

// func (m ReadModel) Resize(width, height int) components.Model {
// 	m.code.Width = width - 1
// 	m.code.Height = height - 1
// 	m.containerWidth = width
// 	m.containerHeight = height
// 	return m
// }

// func (m ReadModel) ContainerWidth() int {
// 	return m.containerWidth
// }

// func (m ReadModel) ContainerHeight() int {
// 	return m.containerHeight
// }

// type Model struct {
// 	containerHeight int
// 	containerWidth  int
// 	theme           theme.Theme
// 	phase           phase_indicator.Model

// 	// ReadModel
// 	readPhase ReadModel
// }

// func NewModel(theme theme.Theme) Model {
// 	return Model{
// 		phase:     phase_indicator.NewModel(theme),
// 		readPhase: NewReadModel(theme),
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	return nil
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {

// 	case messages.IntrospectionEventMsg:
// 		switch evt := msg.Event.(type) {
// 		case compiler_introspection.EventBeginCompileModule:
// 			m.phase.Reset()
// 			m.readPhase.setSourceCode(evt.Origin, evt.SourceCode)
// 		case compiler_introspection.EventEnterPhase:
// 			m.phase.SetCurrentPhase(evt.Phase)
// 		case compiler_introspection.EventEndCompileModule:
// 			m.phase.Finish()
// 		}
// 	}
// 	return m, nil
// }

// func (m Model) View() string {
// 	phaseIndicator := m.phase.View()
// 	// mainView := lipgloss.NewStyle().Height(m.containerHeight).Width(m.containerWidth).Render("Compilation Activity")
// 	mainView := m.readPhase.View()

// 	return lipgloss.JoinVertical(lipgloss.Top, phaseIndicator, mainView)
// }

// func (m Model) Resize(width, height int) activities.Model {
// 	m.phase = m.phase.Resize(width, 1).(phase_indicator.Model)
// 	m.readPhase = m.readPhase.Resize(width/3, height-1).(ReadModel)
// 	m.containerWidth = width
// 	m.containerHeight = height
// 	return m
// }

// func (m Model) ContainerWidth() int {
// 	return m.containerWidth
// }

// func (m Model) ContainerHeight() int {
// 	return m.containerHeight
// }
