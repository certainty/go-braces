package initial_activity

// import (
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/activities"
// 	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/tui/theme"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type Model struct {
// 	containerHeight int
// 	containerWidth  int
// }

// func NewModel(theme theme.Theme) Model {
// 	return Model{
// 		containerWidth:  0,
// 		containerHeight: 0,
// 	}
// }

// func (a Model) Init() tea.Cmd {
// 	return nil
// }

// func (a Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	return a, nil
// }

// func (a Model) View() string {
// 	return lipgloss.NewStyle().
// 		Height(a.containerHeight).
// 		Width(a.containerWidth).
// 		AlignHorizontal(lipgloss.Center).
// 		AlignVertical(lipgloss.Center).
// 		Render("Waiting for next compilation cycle. As soon as the connected compiler begins, the view will change to the corresponding phase introspection.")
// }

// func (a Model) Resize(width, height int) activities.Model {
// 	a.containerWidth = width
// 	a.containerHeight = height
// 	return a
// }

// func (a Model) ContainerWidth() int {
// 	return a.containerWidth
// }

// func (a Model) ContainerHeight() int {
// 	return a.containerHeight
// }
