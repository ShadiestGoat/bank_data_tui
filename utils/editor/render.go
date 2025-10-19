package editor

import (
	"errors"
	"strings"

	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func (c Model) View() string {
	if c.width == 0 {
		return ""
	}

	sections := []string{}
	valid := true

	// last row has special render handling
	for _, row := range c.layout[:len(c.layout) - 1] {
		if len(row) == 1 {
			i := row[0]
			sections = append(sections, renderRowField(c.width, &c.inpFields[i], c.dataFields[i], c.focusedField == i))
		} else {
			parts := make([]string, 0, len(row))
			for _, i := range row {
				parts = append(parts, renderField(&c.inpFields[i], c.dataFields[i], c.focusedField == i))
			}

			sections = append(sections, utils.JoinHorizontalEqualSpread(c.width, parts...))
		}
	}

	btnText := []string{"Save", "Reset"}
	if c.ItemID != "" {
		btnText[0] = "Update"
		btnText[1] = "Delete"
		btnText = append(btnText, "Reset")
	}

	selectedBtn := -1
	if c.focusedField < 0 {
		selectedBtn = -c.focusedField - 1
	}

	sections = append(
		sections,
		scaleButtons(c.width, valid, selectedBtn, btnText),
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

func renderRowField(w int, txt *textinput.Model, data *DataField, selected bool) string {
	fieldStyle := styles.STYLE_FIELD
	if selected {
		fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
	}

	errMsg := ""
	if txt.Err != nil {
		if txt.Value() != "" || errors.Is(txt.Err, ErrRequired) {
			errMsg = txt.Err.Error()
		}
	}

	if data.StyleCB != nil {
		fieldStyle = data.StyleCB(txt.Value(), errMsg, selected, fieldStyle)
	}
	field := fieldStyle.Render(txt.View())

	return utils.JoinHorizontalWithSpacer(
		w, 1,
		data.Title,
		utils.Overflow(
			lipgloss.NewStyle().Faint(true).Italic(true).Bold(errors.Is(txt.Err, APIErr(""))).Render(errMsg),
			w-lipgloss.Width(data.Title)-lipgloss.Width(field)-2,
		)+" ",
		field,
	)
}

func renderField(txt *textinput.Model, data *DataField, selected bool) string {
	fieldStyle := styles.STYLE_FIELD
	if selected {
		fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
	}

	errMsg := ""
	if txt.Err != nil {
		if txt.Value() != "" || errors.Is(txt.Err, ErrRequired) {
			errMsg = txt.Err.Error()
			fieldStyle = fieldStyle.BorderBottom(false)
		}
	}

	if data.StyleCB != nil {
		fieldStyle = data.StyleCB(txt.Value(), errMsg, selected, fieldStyle)
	}

	out := fieldStyle.Render(txt.View())
	if errMsg != "" {
		out += "\n" + utils.Overflow(errMsg, lipgloss.Width(out))
	}

	return out
}
