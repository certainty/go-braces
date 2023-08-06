package colorpalette

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	// frappe colors
	FRAPPE_ROSEWATER = "#f2d5cf"
	FRAPPE_FLAMINGO  = "#eebebe"
	FRAPPE_PINK      = "#f4b8e4"
	FRAPPE_MAUVE     = "#ca9ee6"
	FRAPPE_RED       = "#e78284"
	FRAPPE_MAROON    = "#ea999c"
	FRAPPE_PEACH     = "#ef9f76"
	FRAPPE_YELLOW    = "#e5c890"
	FRAPPE_GREEN     = "#a6d189"
	FRAPPE_TEAL      = "#81c8be"
	FRAPPE_SKY       = "#99d1db"
	FRAPPE_SAPPHIRE  = "#85c1d1"
	FRAPPE_BLUE      = "#8caaee"
	FRAPPE_LAVENDER  = "#babbf1"
	FRAPPE_TEXT      = "#c6d0f5"
	FRAPPE_SUBTEXT1  = "#b5bfe2"
	FRAPPE_SUBTEXT0  = "#a5adce"
	FRAPPE_OVERLAY2  = "#949cbb"
	FRAPPE_OVERLAY1  = "#838ba7"
	FRAPPE_OVERLAY0  = "#737994"
	FRAPPE_SURFACE2  = "#626880"
	FRAPPE_SURFACE1  = "#51576d"
	FRAPPE_SURFACE0  = "#414559"
	FRAPPE_BASE      = "#303446"
	FRAPPE_MANTLE    = "#292c3c"
	FRAPPE_CRUST     = "#232634"

	// machiato colors
	MACCHIATO_ROSEWATER = "#f2d5cf"
	MACCHIATO_FLAMINGO  = "#eebebe"
	MACCHIATO_PINK      = "#f4b8e4"
	MACCHIATO_MAUVE     = "#ca9ee6"
	MACCHIATO_RED       = "#ed8796"
	MACCHIATO_MAROON    = "#ee99a0"
	MACCHIATO_PEACH     = "#f5a97f"
	MACCHIATO_YELLOW    = "#eed49f"
	MACCHIATO_GREEN     = "#a6da95"
	MACCHIATO_TEAL      = "#8bd5ca"
	MACCHIATO_SKY       = "#91d7e3"
	MACCHIATO_SAPPHIRE  = "#7dc4e4"
	MACCHIATO_BLUE      = "#8aadf4"
	MACCHIATO_LAVENDER  = "#b7bdf8"
	MACCHIATO_TEXT      = "#cad3f5"
	MACCHIATO_SUBTEXT1  = "#b8c0e0"
	MACCHIATO_SUBTEXT0  = "#a5adcb"
	MACCHIATO_OVERLAY2  = "#939ab7"
	MACCHIATO_OVERLAY1  = "#8087a2"
	MACCHIATO_OVERLAY0  = "#6e738d"
	MACCHIATO_SURFACE2  = "#5b6078"
	MACCHIATO_SURFACE1  = "#494d64"
	MACCHIATO_SURFACE0  = "#363a4f"
	MACCHIATO_BASE      = "#24273a"
	MACCHIATO_MANTLE    = "#1e2030"
	MACCHIATO_CRUST     = "#181926"
)

type CatpuccinColorScheme struct {
	Rosewater lipgloss.AdaptiveColor
	Flamingo  lipgloss.AdaptiveColor
	Pink      lipgloss.AdaptiveColor
	Mauve     lipgloss.AdaptiveColor
	Red       lipgloss.AdaptiveColor
	Maroon    lipgloss.AdaptiveColor
	Peach     lipgloss.AdaptiveColor
	Yellow    lipgloss.AdaptiveColor
	Green     lipgloss.AdaptiveColor
	Teal      lipgloss.AdaptiveColor
	Sky       lipgloss.AdaptiveColor
	Sapphire  lipgloss.AdaptiveColor
	Blue      lipgloss.AdaptiveColor
	Lavender  lipgloss.AdaptiveColor
	Text      lipgloss.AdaptiveColor
	Subtext1  lipgloss.AdaptiveColor
	Subtext0  lipgloss.AdaptiveColor
	Overlay2  lipgloss.AdaptiveColor
	Overlay1  lipgloss.AdaptiveColor
	Overlay0  lipgloss.AdaptiveColor
	Surface2  lipgloss.AdaptiveColor
	Surface1  lipgloss.AdaptiveColor
	Surface0  lipgloss.AdaptiveColor
	Base      lipgloss.AdaptiveColor
	Mantle    lipgloss.AdaptiveColor
	Crust     lipgloss.AdaptiveColor

	BodyCopy            lipgloss.AdaptiveColor
	MainHeadline        lipgloss.AdaptiveColor
	SubHeadline         lipgloss.AdaptiveColor
	Label               lipgloss.AdaptiveColor
	Subtle              lipgloss.AdaptiveColor
	Link                lipgloss.AdaptiveColor
	Success             lipgloss.AdaptiveColor
	Warning             lipgloss.AdaptiveColor
	Error               lipgloss.AdaptiveColor
	Tag                 lipgloss.AdaptiveColor
	SelectionBackground lipgloss.AdaptiveColor
	Cursor              lipgloss.AdaptiveColor

	CursorText     lipgloss.AdaptiveColor
	ActiveBorder   lipgloss.AdaptiveColor
	InactiveBorder lipgloss.AdaptiveColor
	BellBorder     lipgloss.AdaptiveColor

	Background lipgloss.AdaptiveColor
}

func NewCatpuccinColorScheme() CatpuccinColorScheme {
	scheme := CatpuccinColorScheme{
		Rosewater: lipgloss.AdaptiveColor{Light: FRAPPE_ROSEWATER, Dark: MACCHIATO_ROSEWATER},
		Flamingo:  lipgloss.AdaptiveColor{Light: FRAPPE_FLAMINGO, Dark: MACCHIATO_FLAMINGO},
		Pink:      lipgloss.AdaptiveColor{Light: FRAPPE_PINK, Dark: MACCHIATO_PINK},
		Mauve:     lipgloss.AdaptiveColor{Light: FRAPPE_MAUVE, Dark: MACCHIATO_MAUVE},
		Red:       lipgloss.AdaptiveColor{Light: FRAPPE_RED, Dark: MACCHIATO_RED},
		Maroon:    lipgloss.AdaptiveColor{Light: FRAPPE_MAROON, Dark: MACCHIATO_MAROON},
		Peach:     lipgloss.AdaptiveColor{Light: FRAPPE_PEACH, Dark: MACCHIATO_PEACH},
		Yellow:    lipgloss.AdaptiveColor{Light: FRAPPE_YELLOW, Dark: MACCHIATO_YELLOW},
		Green:     lipgloss.AdaptiveColor{Light: FRAPPE_GREEN, Dark: MACCHIATO_GREEN},
		Teal:      lipgloss.AdaptiveColor{Light: FRAPPE_TEAL, Dark: MACCHIATO_TEAL},
		Sky:       lipgloss.AdaptiveColor{Light: FRAPPE_SKY, Dark: MACCHIATO_SKY},
		Sapphire:  lipgloss.AdaptiveColor{Light: FRAPPE_SAPPHIRE, Dark: MACCHIATO_SAPPHIRE},
		Blue:      lipgloss.AdaptiveColor{Light: FRAPPE_BLUE, Dark: MACCHIATO_BLUE},
		Lavender:  lipgloss.AdaptiveColor{Light: FRAPPE_LAVENDER, Dark: MACCHIATO_LAVENDER},
		Text:      lipgloss.AdaptiveColor{Light: FRAPPE_TEXT, Dark: MACCHIATO_TEXT},
		Subtext1:  lipgloss.AdaptiveColor{Light: FRAPPE_SUBTEXT1, Dark: MACCHIATO_SUBTEXT1},
		Subtext0:  lipgloss.AdaptiveColor{Light: FRAPPE_SUBTEXT0, Dark: MACCHIATO_SUBTEXT0},
		Overlay2:  lipgloss.AdaptiveColor{Light: FRAPPE_OVERLAY2, Dark: MACCHIATO_OVERLAY2},
		Overlay1:  lipgloss.AdaptiveColor{Light: FRAPPE_OVERLAY1, Dark: MACCHIATO_OVERLAY1},
		Overlay0:  lipgloss.AdaptiveColor{Light: FRAPPE_OVERLAY0, Dark: MACCHIATO_OVERLAY0},
		Surface2:  lipgloss.AdaptiveColor{Light: FRAPPE_SURFACE2, Dark: MACCHIATO_SURFACE2},
		Surface1:  lipgloss.AdaptiveColor{Light: FRAPPE_SURFACE1, Dark: MACCHIATO_SURFACE1},
		Surface0:  lipgloss.AdaptiveColor{Light: FRAPPE_SURFACE0, Dark: MACCHIATO_SURFACE0},
		Base:      lipgloss.AdaptiveColor{Light: FRAPPE_BASE, Dark: MACCHIATO_BASE},
		Mantle:    lipgloss.AdaptiveColor{Light: FRAPPE_MANTLE, Dark: MACCHIATO_MANTLE},
		Crust:     lipgloss.AdaptiveColor{Light: FRAPPE_CRUST, Dark: MACCHIATO_CRUST},
	}

	scheme.BodyCopy = scheme.Text
	scheme.MainHeadline = scheme.Text
	scheme.SubHeadline = scheme.Subtext1
	scheme.Label = scheme.Subtext0
	scheme.Subtle = scheme.Overlay1
	scheme.Link = scheme.Blue
	scheme.Success = scheme.Green
	scheme.Warning = scheme.Yellow
	scheme.Error = scheme.Red
	scheme.Tag = scheme.Blue
	scheme.SelectionBackground = scheme.Surface2
	scheme.Cursor = scheme.Rosewater

	scheme.Cursor = scheme.Rosewater
	scheme.CursorText = scheme.Crust
	scheme.ActiveBorder = scheme.Lavender
	scheme.InactiveBorder = scheme.Overlay0
	scheme.BellBorder = scheme.Yellow

	scheme.Background = scheme.Base

	return scheme
}
