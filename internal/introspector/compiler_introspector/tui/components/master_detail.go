package components

import (
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MasterModel struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme
}

type DetailModel struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme
	InputData       string
	OutputData      string
}

type MasterDetailModel struct {
	ContainerWidth  int
	ContainerHeight int
	MasterModel     MasterModel
	DetailModel     DetailModel
	theme           theme.Theme
}

func InitialMainModel(theme theme.Theme) MasterDetailModel {
	return MasterDetailModel{
		theme:       theme,
		MasterModel: InitialMasterModel(theme),
		DetailModel: InitialDetailModel(theme),
	}
}

func (m MasterDetailModel) Init() tea.Cmd {
	return nil
}

func (m MasterDetailModel) Update(msg tea.Msg) (MasterDetailModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	m.MasterModel, cmd = m.MasterModel.Update(msg)
	cmds = append(cmds, cmd)

	m.DetailModel, cmd = m.DetailModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MasterDetailModel) View() string {
	masterWidth := 40 //m.containerWidth - 2
	masterHeight := m.ContainerHeight - 2
	detailWidth := m.ContainerWidth - masterWidth
	detailHeight := m.ContainerHeight - 2

	masterBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, false).BorderForeground(m.theme.Colors.ActiveBorder)
	detailBorder := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, false, true, true).BorderForeground(m.theme.Colors.InactiveBorder)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		masterBorder.Width(masterWidth).Height(masterHeight).Render(m.MasterModel.View()),
		detailBorder.Width(detailWidth).Height(detailHeight).Render(m.DetailModel.View()),
	)

}

func InitialMasterModel(theme theme.Theme) MasterModel {
	return MasterModel{
		theme: theme,
	}
}

func (m MasterModel) Init() tea.Cmd {
	return nil
}

func (m MasterModel) Update(msg tea.Msg) (MasterModel, tea.Cmd) {
	return m, nil
}

func (m MasterModel) View() string {
	return "phases"
}

func InitialDetailModel(theme theme.Theme) DetailModel {
	return DetailModel{
		theme:      theme,
		InputData:  "hello",
		OutputData: "world",
	}
}

func (m DetailModel) Init() tea.Cmd {
	return nil
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd) {
	return m, nil
}

func (m DetailModel) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.InputData,
		m.OutputData,
	)
}
