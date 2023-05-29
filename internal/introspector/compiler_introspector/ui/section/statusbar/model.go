package statusbar

import (
	"github.com/certainty/go-braces/internal/introspection/compiler_introspection"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/philistino/teacup/statusbar"
)

type Model struct {
	containerWidth  int
	containerHeight int
	theme           theme.Theme

	Phase         compiler_introspection.CompilationPhase
	Mode          common.Mode
	globalKeyMap  common.KeyMap
	contextKeyMap common.KeyMap
	RequestState  common.RequestStatus

	IsConnected    bool
	err            error
	RequestSpinner spinner.Model
	impl           statusbar.Bubble
}

func New(theme theme.Theme, globalKeyMap common.KeyMap, contextKeyMap common.KeyMap) Model {
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

	return Model{
		theme:          theme,
		Phase:          "",
		Mode:           common.WaitingMode,
		err:            nil,
		IsConnected:    false,
		globalKeyMap:   globalKeyMap,
		contextKeyMap:  contextKeyMap,
		RequestSpinner: requestSpinner,
		impl:           sb,
	}
}
