package tui

import (
	"log"

	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/activities"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/activities/initial_activity"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/commands"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components/header"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components/statusbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/messages"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	shared_components "github.com/certainty/go-braces/internal/introspector/tui_shared/components"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type IntrospectionKeyMap struct {
	ToggleSingleStepping key.Binding
	Step                 key.Binding
	Continue             key.Binding
	Reset                key.Binding
}

type GlobalKeyMap struct {
	Quit           key.Binding
	Help           key.Binding
	Documentation  key.Binding
	ToggleEventLog key.Binding
}

type ActivityName int

const (
	InitialActivityName ActivityName = iota
)

type TUIModel struct {
	// tui state
	width               int
	height              int
	theme               theme.Theme
	globalKeyMap        GlobalKeyMap
	introspectionKeyMap IntrospectionKeyMap

	// components
	headerModel    header.Model
	infoModel      components.InfoModel
	statusBarModel statusbar.Model
	eventLogModel  shared_components.EventLogModel

	// activities
	activities          []activities.Model
	currentActivityName ActivityName

	eventLogVisible bool

	// data
	client                 *compiler_introspection.Client
	currentCompilerOptions []string
	currentInput           *input.Input
	currentPhase           string
	requestState           components.IntrospectionRequestState
	singleStepMode         bool
	connected              bool
}

// Messages
// Commands

func InitialTUIModel(client *compiler_introspection.Client) TUIModel {
	theme := theme.NewCatpuccinTheme()
	currentInput := input.NewReplInput(1, "#t")
	currentCompilerOptions := []string{}
	connected := client.IsConnected()

	globalKeyMap := GlobalKeyMap{
		Quit:           key.NewBinding(key.WithKeys("Q"), key.WithHelp("[Q]uit", "Quit the program")),
		Help:           key.NewBinding(key.WithKeys("H"), key.WithHelp("[H]elp", "Open help dialog")),
		Documentation:  key.NewBinding(key.WithKeys("D"), key.WithHelp("[D]ocumentation", "Open documentation dialog")),
		ToggleEventLog: key.NewBinding(key.WithKeys("E"), key.WithHelp("[E]vents", "Toggle event log")),
	}

	introspectionKeyMap := IntrospectionKeyMap{
		ToggleSingleStepping: key.NewBinding(key.WithKeys("S"), key.WithHelp("[S]ingle", "Toggle single stepping")),
		Step:                 key.NewBinding(key.WithKeys("N"), key.WithHelp("[N]ext", "Step to next phase")),
		Continue:             key.NewBinding(key.WithKeys("C"), key.WithHelp("[C]ontinue", "Continue to end")),
		Reset:                key.NewBinding(key.WithKeys("R"), key.WithHelp("[R]eset", "Reset introspection")),
	}
	introspectionKeyMap.Step.SetEnabled(false)
	introspectionKeyMap.Continue.SetEnabled(true)
	introspectionKeyMap.Reset.SetEnabled(false)

	shortcuts := []*key.Binding{
		&introspectionKeyMap.ToggleSingleStepping,
		&introspectionKeyMap.Step,
		&introspectionKeyMap.Continue,
		&introspectionKeyMap.Reset,
		&globalKeyMap.ToggleEventLog,
		&globalKeyMap.Quit,
		&globalKeyMap.Help,
	}

	activities := []activities.Model{
		initial_activity.NewModel(theme),
	}

	return TUIModel{
		width:               10,
		height:              80,
		theme:               theme,
		globalKeyMap:        globalKeyMap,
		introspectionKeyMap: introspectionKeyMap,
		headerModel:         header.NewModel(theme, "(Go-Braces-Introspect 'Compiler)", connected),
		infoModel:           components.InitialInfoModel(theme, currentInput, currentCompilerOptions),
		statusBarModel:      statusbar.NewModel(theme, shortcuts),
		eventLogModel:       shared_components.InitialEventLogModel(),

		activities:          activities,
		currentActivityName: InitialActivityName,

		eventLogVisible:        false,
		client:                 client,
		currentCompilerOptions: currentCompilerOptions,
		singleStepMode:         false,
		requestState:           components.NoRequest,
		currentInput:           currentInput,
		currentPhase:           "waiting",
	}
}

func (m TUIModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, commands.DoGetEvent(m.client))
	cmds = append(cmds, commands.DoTick())

	for i := range m.activities {
		cmds = append(cmds, m.activities[i].Init())
	}

	// Do I need to do this?
	cmd := m.statusBarModel.Init()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.client.IsConnected() {
		m.connected = true
	} else {
		m.connected = false
		m.statusBarModel.Errors = "Lost connection to introspection server ..."
	}
	m.headerModel.Connected = m.connected

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.propagateSizeChange()
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.globalKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.globalKeyMap.ToggleEventLog):
			m.eventLogVisible = !m.eventLogVisible
			m.propagateSizeChange()
		case key.Matches(msg, m.introspectionKeyMap.Continue):
			m.singleStepMode = false
			return m, commands.DoBreakpointContinue(m.client)
		default:
			m, cmds := m.propagateUpdate(msg)
			return m, tea.Batch(cmds...)
		}
	case messages.TickMsg:
		return m, tea.Batch(commands.DoTick())

	case messages.IntrospectionConnectionEstablished:
		m.connected = true

	case messages.IntrospectionConnectionError:
		m.connected = false
		m.statusBarModel.Errors = msg.Err.Error()

	case messages.IntrospectionEventMsg:
		m.eventLogModel, _ = m.eventLogModel.Update(shared_components.NewEventMessage(msg))
		cmds = append(cmds, commands.DoGetEvent(m.client))

		m, propagatedCommands := m.propagateUpdate(msg)
		cmds = append(cmds, propagatedCommands...)

		return m, tea.Batch(cmds...)
	default:
		m, propagatedCommands := m.propagateUpdate(msg)
		cmds = append(cmds, propagatedCommands...)
		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m TUIModel) propagateUpdate(msg tea.Msg) (tea.Model, []tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	newHeaderModel, cmd := m.headerModel.Update(msg)
	m.headerModel = newHeaderModel.(header.Model)
	cmds = append(cmds, cmd)

	m.infoModel, cmd = m.infoModel.Update(msg)
	cmds = append(cmds, cmd)

	// update current activity
	activityModel, cmd := m.activities[m.currentActivityName].Update(msg)
	m.activities[m.currentActivityName] = activityModel.(activities.Model)
	cmds = append(cmds, cmd)

	newStatusbarModel, cmd := m.statusBarModel.Update(msg)
	m.statusBarModel = newStatusbarModel.(statusbar.Model)
	cmds = append(cmds, cmd)

	return m, cmds
}

func (m *TUIModel) currentActivity() activities.Model {
	return m.activities[m.currentActivityName]
}

func (m *TUIModel) propagateSizeChange() {
	width := m.width
	height := m.height

	m.headerModel.Resize(width, 1)

	m.infoModel.ContainerWidth = width
	m.infoModel.ContainerHeight = 3

	m.statusBarModel = m.statusBarModel.Resize(width, 1).(statusbar.Model)

	log.Printf("Event log visible: %v\n", m.eventLogVisible)
	if m.eventLogVisible {
		m.eventLogModel.ContainerHeight = 20
	} else {
		m.eventLogModel.ContainerHeight = 0
	}
	m.eventLogModel.ContainerWidth = width

	// all but the heights above
	activityContainerWidth := width - 2
	activityContainerHeight := height - m.headerModel.ContainerHeight() - m.infoModel.ContainerHeight - m.statusBarModel.ContainerHeight() - m.eventLogModel.ContainerHeight
	m.activities[m.currentActivityName] = m.currentActivity().Resize(activityContainerWidth, activityContainerHeight)
}

func (m TUIModel) View() string {
	termWidth := m.width
	termHeight := m.height

	headerView := m.headerModel.View()
	infoView := m.infoModel.View()
	statusBarView := m.statusBarModel.View()

	activityView := m.currentActivity().View()

	components := []string{headerView, infoView, activityView}
	if m.eventLogVisible {
		components = append(components, m.eventLogModel.View())
	}
	components = append(components, statusBarView)

	// Assemble the components
	content := lipgloss.NewStyle().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, components...))

	return content
}
