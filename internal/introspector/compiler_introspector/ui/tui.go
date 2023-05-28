package ui

// import (
// 	"github.com/certainty/go-braces/internal/compiler/input"
// 	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/activities"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/activities/compile_activity"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/activities/initial_activity"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/commands"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/components"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/components/compilation_info"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/components/header"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/messages"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/eventlog"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
// 	"github.com/charmbracelet/bubbles/key"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type IntrospectionKeyMap struct {
// 	ToggleSingleStepping key.Binding
// 	Step                 key.Binding
// 	Continue             key.Binding
// 	Reset                key.Binding
// }

// type GlobalKeyMap struct {
// 	Quit           key.Binding
// 	Help           key.Binding
// 	Documentation  key.Binding
// 	ToggleEventLog key.Binding
// }

// type ActivityName int

// const (
// 	InitialActivity ActivityName = iota
// 	CompileActivity
// )

// type Model struct {
// 	// tui state
// 	width               int
// 	height              int
// 	theme               theme.Theme
// 	globalKeyMap        GlobalKeyMap
// 	introspectionKeyMap IntrospectionKeyMap

// 	// components
// 	headerModel   header.Model
// 	infoModel     compilation_info.Model
// 	eventLogModel eventlog.Model

// 	// activities
// 	activities          []activities.Model
// 	currentActivityName ActivityName

// 	eventLogVisible bool

// 	// data
// 	client                 *compiler_introspection.Client
// 	currentCompilerOptions []string
// 	currentInput           *input.Input
// 	currentPhase           string
// 	requestState           components.IntrospectionRequestState
// 	singleStepMode         bool
// 	connected              bool
// }

// // Messages
// // Commands

// func NewModel(client *compiler_introspection.Client) Model {
// 	theme := theme.NewCatpuccinTheme()
// 	currentInput := input.NewReplInput(1, "#t")
// 	currentCompilerOptions := []string{}
// 	connected := client.IsConnected()

// 	globalKeyMap := GlobalKeyMap{
// 		Quit:           key.NewBinding(key.WithKeys("Q"), key.WithHelp("[Q]uit", "Quit the program")),
// 		Help:           key.NewBinding(key.WithKeys("H"), key.WithHelp("[H]elp", "Open help dialog")),
// 		Documentation:  key.NewBinding(key.WithKeys("D"), key.WithHelp("[D]ocumentation", "Open documentation dialog")),
// 		ToggleEventLog: key.NewBinding(key.WithKeys("E"), key.WithHelp("[E]vents", "Toggle event log")),
// 	}

// 	introspectionKeyMap := IntrospectionKeyMap{
// 		ToggleSingleStepping: key.NewBinding(key.WithKeys("S"), key.WithHelp("[S]ingle", "Toggle single stepping")),
// 		Step:                 key.NewBinding(key.WithKeys("N"), key.WithHelp("[N]ext", "Step to next phase")),
// 		Continue:             key.NewBinding(key.WithKeys("C"), key.WithHelp("[C]ontinue", "Continue to end")),
// 		Reset:                key.NewBinding(key.WithKeys("R"), key.WithHelp("[R]eset", "Reset introspection")),
// 	}
// 	introspectionKeyMap.Step.SetEnabled(false)
// 	introspectionKeyMap.Continue.SetEnabled(true)
// 	introspectionKeyMap.Reset.SetEnabled(false)

// 	activities := []activities.Model{
// 		initial_activity.NewModel(theme),
// 		compile_activity.NewModel(theme),
// 	}

// 	return Model{
// 		width:                  10,
// 		height:                 80,
// 		theme:                  theme,
// 		globalKeyMap:           globalKeyMap,
// 		introspectionKeyMap:    introspectionKeyMap,
// 		headerModel:            header.NewModel(theme, "(Go-Braces-Introspect 'Compiler)", connected),
// 		infoModel:              compilation_info.NewModel(theme, currentInput, currentCompilerOptions),
// 		eventLogModel:          eventlog.NewModel(),
// 		activities:             activities,
// 		currentActivityName:    InitialActivity,
// 		eventLogVisible:        false,
// 		client:                 client,
// 		currentCompilerOptions: currentCompilerOptions,
// 		currentInput:           currentInput,
// 		currentPhase:           "waiting",
// 		requestState:           components.NoRequest,
// 		singleStepMode:         false,
// 		connected:              connected,
// 	}
// }

// func (m Model) Init() tea.Cmd {
// 	var cmds []tea.Cmd
// 	cmds = append(cmds, commands.DoGetEvent(m.client))
// 	cmds = append(cmds, commands.DoTick())

// 	for i := range m.activities {
// 		cmds = append(cmds, m.activities[i].Init())
// 	}

// 	return tea.Batch(cmds...)
// }

// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd

// 	switch msg := msg.(type) {

// 	case tea.WindowSizeMsg:
// 		m.width = msg.Width
// 		m.height = msg.Height
// 		m.propagateSizeChange()
// 		return m, tea.Batch(cmds...)

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, m.globalKeyMap.Quit):
// 			return m, tea.Quit
// 		case key.Matches(msg, m.globalKeyMap.ToggleEventLog):
// 			m.eventLogVisible = !m.eventLogVisible
// 			m.propagateSizeChange()
// 		case key.Matches(msg, m.introspectionKeyMap.Continue):
// 			m.singleStepMode = false
// 			return m, commands.DoBreakpointContinue(m.client)
// 		default:
// 			m, cmds := m.propagateUpdate(msg)
// 			return m, tea.Batch(cmds...)
// 		}
// 	case messages.TickMsg:
// 		return m, tea.Batch(commands.DoTick())

// 	case messages.IntrospectionConnectionEstablished:
// 		m.connected = true

// 	case messages.IntrospectionConnectionError:
// 		m.connected = false

// 	case messages.IntrospectionEventMsg:
// 		newEventlogModel, _ := m.eventLogModel.Update(eventlog.NewEventMessage(msg))
// 		m.eventLogModel = newEventlogModel.(eventlog.Model)

// 		switch msg.Event.(type) {
// 		case compiler_introspection.EventBeginCompileModule:
// 			if m.currentActivityName != CompileActivity {
// 				m.switchActivity(CompileActivity)
// 			}
// 		}
// 		cmds = append(cmds, commands.DoGetEvent(m.client))
// 		m, propagatedCommands := m.propagateUpdate(msg)
// 		cmds = append(cmds, propagatedCommands...)

// 		return m, tea.Batch(cmds...)
// 	default:
// 		m, propagatedCommands := m.propagateUpdate(msg)
// 		cmds = append(cmds, propagatedCommands...)
// 		return m, tea.Batch(cmds...)
// 	}

// 	return m, tea.Batch(cmds...)
// }

// func (m Model) propagateUpdate(msg tea.Msg) (tea.Model, []tea.Cmd) {
// 	var cmds []tea.Cmd
// 	var cmd tea.Cmd

// 	newHeaderModel, cmd := m.headerModel.Update(msg)
// 	m.headerModel = newHeaderModel.(header.Model)
// 	cmds = append(cmds, cmd)

// 	newInfoModel, cmd := m.infoModel.Update(msg)
// 	m.infoModel = newInfoModel.(compilation_info.Model)
// 	cmds = append(cmds, cmd)

// 	// update current activity
// 	activityModel, cmd := m.activities[m.currentActivityName].Update(msg)
// 	m.activities[m.currentActivityName] = activityModel.(activities.Model)
// 	cmds = append(cmds, cmd)

// 	return m, cmds
// }

// func (m *Model) currentActivity() activities.Model {
// 	return m.activities[m.currentActivityName]
// }

// func (m *Model) switchActivity(name ActivityName) {
// 	m.currentActivityName = name
// 	m.propagateSizeChange()
// }

// func (m *Model) propagateSizeChange() {
// 	width := m.width
// 	height := m.height

// 	m.headerModel = m.headerModel.Resize(width, 1).(header.Model)
// 	m.infoModel = m.infoModel.Resize(width, 3).(compilation_info.Model)

// 	if m.eventLogVisible {
// 		m.eventLogModel = m.eventLogModel.Resize(width, 20).(eventlog.Model)
// 	} else {
// 		m.eventLogModel = m.eventLogModel.Resize(width, 0).(eventlog.Model)
// 	}

// 	// all but the heights above
// 	activityContainerWidth := width - 2
// 	activityContainerHeight := height - m.headerModel.ContainerHeight() - m.infoModel.ContainerHeight() - m.eventLogModel.ContainerHeight()
// 	m.activities[m.currentActivityName] = m.currentActivity().Resize(activityContainerWidth, activityContainerHeight)
// }

// func (m Model) View() string {
// 	termWidth := m.width
// 	termHeight := m.height

// 	headerView := m.headerModel.View()
// 	infoView := m.infoModel.View()

// 	activityView := m.currentActivity().View()

// 	components := []string{headerView, infoView, activityView}
// 	if m.eventLogVisible {
// 		components = append(components, m.eventLogModel.View())
// 	}

// 	// Assemble the components
// 	content := lipgloss.NewStyle().
// 		Width(termWidth).
// 		Height(termHeight).
// 		Render(lipgloss.JoinVertical(lipgloss.Top, components...))

// 	return content
// }