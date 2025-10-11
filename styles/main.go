package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	COLOR_MAIN = lipgloss.Color("#6557f9")
	COLOR_WRONG = lipgloss.ANSIColor(1)
	COLOR_DISABLED = lipgloss.ANSIColor(8)
)

var (
	STYLE_FIELD = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.DoubleBorder()).Background(lipgloss.NoColor{})
	STYLE_BTN   = STYLE_FIELD.Padding(1, 2)

	style_base_btn_selected = STYLE_BTN.Foreground(lipgloss.NoColor{})
	STYLE_BTN_SELECTED = style_base_btn_selected.Background(COLOR_MAIN)
	STYLE_BTN_SELECTED_DISABLED = style_base_btn_selected.Background(COLOR_DISABLED).BorderForeground(COLOR_DISABLED)
	STYLE_BTN_SELECTED_BAD = style_base_btn_selected.Background(COLOR_WRONG)
)
