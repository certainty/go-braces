package theme

import (
	"github.com/certainty/go-braces/internal/introspector/tui_shared/colorpalette"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Colors     colorpalette.CatpuccinColorScheme
	Background lipgloss.Style
	Foreground lipgloss.Style
	Accent     lipgloss.Style
	Highlight  lipgloss.Style

	Label lipgloss.Style

	// specific styles
	Header lipgloss.Style
	Title  lipgloss.Style

	// statusBar
	Statusbar               lipgloss.Style
	StatusbarShortcutPrefix lipgloss.Style

	// info
	Info lipgloss.Style
}

// A theme combines colors and font styles
func NewCatpuccinTheme() Theme {
	colors := colorpalette.NewCatpuccinColorScheme()
	header := lipgloss.NewStyle().Background(colors.Crust).Foreground(colors.MainHeadline)
	title := header.Copy().PaddingLeft(1).PaddingRight(2).Bold(true)
	statusBar := lipgloss.NewStyle().Background(colors.Background).Foreground(colors.Text)
	label := lipgloss.NewStyle().Foreground(colors.Text).Bold(true)

	return Theme{
		Colors:     colors,
		Background: lipgloss.NewStyle().Background(colors.Base),
		Foreground: lipgloss.NewStyle().Foreground(colors.Text),
		Accent:     lipgloss.NewStyle().Foreground(colors.Pink),

		Label: label,

		// concrete styles
		Header:                  header,
		Title:                   title,
		Info:                    lipgloss.NewStyle().Foreground(colors.Subtext1),
		Statusbar:               statusBar,
		StatusbarShortcutPrefix: statusBar.Copy().Foreground(colors.Red).Bold(true),
	}
}
