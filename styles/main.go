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
	STYLE_BTN_DISABLED = STYLE_BTN.Foreground(COLOR_DISABLED).BorderForeground(COLOR_DISABLED)
	STYLE_BTN_SELECTED = style_base_btn_selected.Background(COLOR_MAIN)
	STYLE_BTN_SELECTED_DISABLED = style_base_btn_selected.Background(COLOR_DISABLED).BorderForeground(COLOR_DISABLED)
	STYLE_BTN_SELECTED_BAD = style_base_btn_selected.Background(COLOR_WRONG)
)

func StyleBtn(disabled, selected, bad, small bool) lipgloss.Style {
	style := STYLE_FIELD
	if !small {
		style = STYLE_BTN
	}

	if disabled && !selected {
		style = style.Foreground(COLOR_DISABLED)
	} else if selected {
		style = style.Foreground(lipgloss.NoColor{})
	}

	var color lipgloss.TerminalColor = COLOR_MAIN
	switch {
	case disabled:
		color = COLOR_DISABLED
	case bad:
		color = COLOR_WRONG
	}

	if selected {
		style = style.Background(color)
	}

	return style.BorderForeground(color)
}
