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
	for _, row := range c.layout[:len(c.layout)-1] {
		if len(row) == 1 {
			i := row[0]
			txt := &c.inpFields[i]
			sections = append(sections, renderRowField(c.width, txt, c.dataFields[i], c.focusedField == i))
			if txt.Err != nil {
				valid = false
			}
		} else {
			parts := make([]string, 0, len(row))
			for _, i := range row {
				txt := &c.inpFields[i]
				parts = append(parts, renderField(txt, c.dataFields[i], c.focusedField == i))
				if txt.Err != nil {
					valid = false
				}
			}

			sections = append(sections, utils.JoinHorizontalEqualSpread(c.width, parts...))
		}
	}

	btnText := []string{"Save", "Reset"}
	btnIDs := []int{BTN_SAVE, BTN_RESET}
	if c.ItemID != "" {
		btnText[0] = "Update"
		btnText[1] = "Delete"
		btnText = append(btnText, "Reset")
		btnIDs[1] = BTN_DEL
		btnIDs = append(btnIDs, BTN_RESET)
	}

	selectedBtn := -1
	for i, id := range btnIDs {
		if id == c.focusedField {
			selectedBtn = i
			break
		}
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

	err := renderErr(txt)
	if data.StyleCB != nil {
		fieldStyle = data.StyleCB(txt.Value(), txt.Err, selected, fieldStyle)
	}
	field := fieldStyle.Render(renderTextField(txt))

	return utils.JoinHorizontalWithSpacer(
		w, 1,
		data.Title,
		utils.Overflow(
			err,
			w-lipgloss.Width(data.Title)-lipgloss.Width(field)-2,
		)+" ",
		field,
	)
}

func renderErr(txt *textinput.Model) string {
	if txt.Err == nil {
		return ""
	}
	if txt.Value() == "" && !errors.Is(txt.Err, ErrRequired{}) {
		return ""
	}

	return lipgloss.NewStyle().Foreground(styles.COLOR_WRONG).Faint(true).Italic(true).Bold(
		errors.Is(txt.Err, APIErr("")),
	).Render(txt.Err.Error())
}

func renderField(txt *textinput.Model, data *DataField, selected bool) string {
	fieldStyle := styles.STYLE_FIELD
	if selected {
		fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
	}
	if data.StyleCB != nil {
		fieldStyle = data.StyleCB(txt.Value(), txt.Err, selected, fieldStyle)
	}
	
	err := renderErr(txt)
	if err != "" {
		fieldStyle = fieldStyle.BorderBottom(false)
	}
	if txt.Value() != "" {
		fieldStyle = fieldStyle.BorderTop(false)
	}
	
	out := fieldStyle.Render(renderTextField(txt))
	if err != "" {
		out += "\n" + fakeBorder(false, fieldStyle, renderErr(txt), lipgloss.Width(out))
	}
	if txt.Value() == "" {
		return out
	}

	titleStyle := lipgloss.NewStyle().Faint(true)
	if selected {
		titleStyle = titleStyle.Foreground(styles.COLOR_MAIN)
	}

	return fakeBorder(true, fieldStyle, titleStyle.Render(data.Title), lipgloss.Width(out)) + "\n" + out
}

func renderTextField(txt *textinput.Model) string {
	// Fuck your text field and it's horse
	return lipgloss.NewStyle().Width(txt.Width + 1).MaxWidth(txt.Width + 1).Inline(true).Render(txt.View())
}

func extraFieldLength(txt *textinput.Model, data *DataField) int {
	fieldStyle := styles.STYLE_FIELD
	if data.StyleCB != nil {
		fieldStyle = data.StyleCB(txt.Value(), txt.Err, false, fieldStyle)
	}

	return lipgloss.Width(fieldStyle.Render(""))
}

func fakeBorder(top bool, style lipgloss.Style, str string, totalWidth int) string {
	var (
		horz       string
		leftCorner  string
		rightCorner string
	)

	bo := style.GetBorderStyle()
	if top {
		horz = bo.Top
		leftCorner = bo.TopLeft
		rightCorner = bo.TopRight
	} else {
		horz = bo.Bottom
		leftCorner = bo.BottomLeft
		rightCorner = bo.BottomRight
	}

	newStyle := lipgloss.NewStyle().Foreground(style.GetBorderLeftForeground())

	return newStyle.Render(leftCorner+horz) + " " +
		lipgloss.NewStyle().Width(totalWidth - 6).Render(utils.Overflow(str, totalWidth-6)) +
		" " + newStyle.Render(horz+rightCorner)
}
