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

	for i, f := range c.dataFields {
		fieldStyle := styles.STYLE_FIELD
		txt := c.inpFields[i]
		if c.focusedField == i {
			fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
		}

		errMsg := ""
		apiErr := false
		if txt.Err != nil {
			valid = false

			if txt.Value() != "" || errors.Is(txt.Err, ErrRequired) {
				errMsg = txt.Err.Error()
			}
		} else if f.apiErr != "" {
			apiErr = true
			errMsg = f.apiErr
		}

		field := fieldStyle.Render(txt.View())

		sections = append(sections, utils.JoinHorizontalSpread(
			c.width, 1,
			f.Title,
			utils.Overflow(
				lipgloss.NewStyle().Faint(true).Italic(true).Bold(apiErr).Render(errMsg),
				c.width - lipgloss.Width(f.Title) - lipgloss.Width(field) - 2,
			) + " ",
			field,
		))
	}

	btns := []string{"Save", "Reset"}
	if c.id != nil {
		btns[0] = "Update"
	}

	for i, b := range btns {
		fieldStyle := styles.StyleBtn(
			!valid,
			i == c.focusedField-len(c.dataFields),
		)

		btns[i] = fieldStyle.Render(b)
	}

	sections = append(
		sections,
		utils.JoinHorizontal2(c.width, btns[0], btns[1]),
	)

	return strings.Join(sections, "\n\n")
}
