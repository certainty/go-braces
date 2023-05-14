package tui

import (
	"github.com/certainty/go-braces/internal/compiler/input"
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/commands"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/components"
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

type TUIModel struct {
	// tui state
	width               int
	height              int
	theme               theme.Theme
	globalKeyMap        GlobalKeyMap
	introspectionKeyMap IntrospectionKeyMap

	// components
	headerModel              components.HeaderModel
	infoModel                components.InfoModel
	mainModel                components.MasterDetailModel
	statusBarModel           components.StatusBarModel
	helpDialogModel          components.HelpDialogModel
	documentationDialogModel components.DocumentationDialogModel
	eventLogModel            shared_components.EventLogModel

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

	return TUIModel{
		width:                    10,
		height:                   80,
		theme:                    theme,
		globalKeyMap:             globalKeyMap,
		introspectionKeyMap:      introspectionKeyMap,
		headerModel:              components.InitialHeaderModel(theme, "(Go-Braces-Introspect 'Compiler)", connected),
		infoModel:                components.InitialInfoModel(theme, currentInput, currentCompilerOptions),
		mainModel:                components.InitialMainModel(theme),
		statusBarModel:           components.InitialStatusBarModel(theme, shortcuts),
		helpDialogModel:          components.InitialHelpDialogModel(theme),
		documentationDialogModel: components.InitialDocumentationDialogModel(theme),
		eventLogModel:            shared_components.InitialEventLogModel(),

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
		case key.Matches(msg, m.globalKeyMap.Help):
			m.helpDialogModel.Active = !m.helpDialogModel.Active
		case key.Matches(msg, m.globalKeyMap.Documentation):
			m.documentationDialogModel.Active = !m.documentationDialogModel.Active
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
	m.headerModel, cmd = m.headerModel.Update(msg)
	cmds = append(cmds, cmd)

	m.infoModel, cmd = m.infoModel.Update(msg)
	cmds = append(cmds, cmd)

	m.mainModel, cmd = m.mainModel.Update(msg)
	cmds = append(cmds, cmd)

	m.statusBarModel, cmd = m.statusBarModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, cmds
}

func (m *TUIModel) propagateSizeChange() {
	width := m.width
	height := m.height

	m.headerModel.ContainerWidth = width
	m.headerModel.ContainerHeight = 1

	m.infoModel.ContainerWidth = width
	m.infoModel.ContainerHeight = 3

	m.statusBarModel.ContainerWidth = width
	m.statusBarModel.ContainerHeight = 1

	if m.eventLogVisible {
		m.eventLogModel.ContainerHeight = 20
	} else {
		m.eventLogModel.ContainerHeight = 0
	}
	m.eventLogModel.ContainerWidth = width

	// all but the heights above
	m.mainModel.ContainerWidth = width - 2
	m.mainModel.ContainerHeight = height - m.headerModel.ContainerHeight - m.infoModel.ContainerHeight - m.statusBarModel.ContainerHeight - m.eventLogModel.ContainerHeight

	m.helpDialogModel.ContainerWidth = width
	m.helpDialogModel.ContainerHeight = height

	m.documentationDialogModel.ContainerWidth = width
	m.documentationDialogModel.ContainerHeight = height

}

func (m TUIModel) View() string {
	termWidth := m.width
	termHeight := m.height

	headerView := m.headerModel.View()
	infoView := m.infoModel.View()
	statusBarView := m.statusBarModel.View()
	mainView := m.mainModel.View()

	components := []string{headerView, infoView, mainView}
	if m.eventLogVisible {
		components = append(components, m.eventLogModel.View())
	}
	components = append(components, statusBarView)

	// Assemble the components
	content := lipgloss.NewStyle().
		Width(termWidth).
		Height(termHeight).
		Render(lipgloss.JoinVertical(lipgloss.Top, components...))

	// Render help and documentation dialogs if they are active
	// if m.helpDialogModel.active {
	// 	content = lipgloss.Stack(content, m.helpDialogModel.View())
	// }
	// if m.documentationDialogModel.active {
	// 	content = lipgloss.Stack(content, m.documentationDialogModel.View())
	// }

	return content
}
