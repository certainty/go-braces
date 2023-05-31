package topbar_test

import (
	"testing"

	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/common"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/section/topbar"
	"github.com/certainty/go-braces/internal/introspector/compiler_introspector/ui/theme"
	"github.com/gkampitakis/go-snaps/snaps"
)

func TestModel(t *testing.T) {
	theme := theme.NewCatpuccinTheme()
	model := topbar.New(theme, "My Title")

	snaps.MatchSnapshot(t, model.View())
}

func TestModelWithConnectedClient(t *testing.T) {
	theme := theme.NewCatpuccinTheme()
	model := topbar.New(theme, "My Title")
	model.Update(common.MsgClientConnected(true))

	snaps.MatchSnapshot(t, model.View())
}
