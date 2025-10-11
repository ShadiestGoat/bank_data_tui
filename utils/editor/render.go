package editor

import (
	"errors"
	"strings"

	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/lipgloss"
)

func (c Model) View() string {
	if c.width == 0 {
		return ""
	}

	sections := []string{}
	valid := true

	for i, t := range c.titles {
		fieldStyle := styles.STYLE_FIELD
		txt := c.textFields[i]
		if c.focusedField == i {
			fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
		}

		errMsg := ""
		if txt.Err != nil {
			valid = false

			if txt.Value() != "" || errors.Is(txt.Err, ErrRequired) {
				errMsg = txt.Err.Error() + " "
			}
		}

		sections = append(sections, utils.JoinHorizontalSpread(
			c.width, 1,
			t,
			lipgloss.NewStyle().Faint(true).Italic(true).Render(errMsg),
			fieldStyle.Render(txt.View()),
		))
	}

	btns := []string{"Save", "Reset"}
	if c.id != nil {
		btns[0] = "Update"
	}

	for i, b := range btns {
		fieldStyle := styles.STYLE_BTN
		if i == c.focusedField-len(c.editableFields) {
			if i == 0 {
				if valid {
					fieldStyle = styles.STYLE_BTN_SELECTED
				} else {
					fieldStyle = styles.STYLE_BTN_SELECTED_DISABLED
				}
			} else {
				fieldStyle = styles.STYLE_BTN_SELECTED_BAD
			}
		}

		btns[i] = fieldStyle.Render(b)
	}

	sections = append(
		sections,
		utils.JoinHorizontal2(c.width, btns[0], btns[1]),
	)

	return strings.Join(sections, "\n\n")
}
