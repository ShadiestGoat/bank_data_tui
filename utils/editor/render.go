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

		sections = append(sections, utils.JoinHorizontalWithSpacer(
			c.width, 1,
			f.Title,
			utils.Overflow(
				lipgloss.NewStyle().Faint(true).Italic(true).Bold(apiErr).Render(errMsg),
				c.width-lipgloss.Width(f.Title)-lipgloss.Width(field)-2,
			)+" ",
			field,
		))
	}

	btnText := []string{"Save", "Reset"}
	if c.ItemID != "" {
		btnText[0] = "Update"
		btnText[1] = "Delete"
		btnText = append(btnText, "Reset")
	}

	sections = append(
		sections,
		scaleButtons(c.width, valid, c.focusedField - len(c.dataFields), btnText),
	)

	return strings.Join(sections, "\n\n")
}

func scaleButtons(w int, valid bool, selectedBtn int, btnText []string) string {
	if t := renderButtons(w, valid, selectedBtn, false, btnText); t != "" {
		return t
	}

	return renderButtons(w, valid, selectedBtn, true, btnText)
}

func renderButtons(w int, valid bool, selectedBtn int, small bool, btnText []string) string {
	btns := make([]string, len(btnText))
	for i, t := range btnText {
		btns[i] = styles.StyleBtn(
			!valid && i == 0,
			i == selectedBtn,
			i != 0,
			small,
		).Render(t)
	}

	return utils.JoinHorizontalEqualSpread(w, btns...)
}