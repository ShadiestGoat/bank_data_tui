package login

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
)

var (
	// STYLE_MOD_OK       = lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	STYLE_MOD_DISABLED = lipgloss.NewStyle().BorderForeground(lipgloss.ANSIColor(8)).Background(lipgloss.ANSIColor(8))
	STYLE_MOD_WRONG    = lipgloss.NewStyle().BorderForeground(lipgloss.ANSIColor(9)).Background(lipgloss.ANSIColor(9))
)

var _ utils.Screen = &Model{} // compile check

type Model struct {
	focusedField int
	// 0 = ok
	// 1 = loading
	// 2 = wrong
	state int

	inpName textinput.Model
	inpPass textinput.Model
}

func NewScreenLogin() *Model {
	inpName := textinput.New()
	inpPass := textinput.New()

	for i, inp := range []*textinput.Model{&inpName, &inpPass} {
		inp.SetWidth(15)
		inp.Prompt = ""
		inp.SetVirtualCursor(false)
		inp.SetStyles(textinput.Styles{
			Focused: textinput.StyleState{
				Text:        styles.S_TEXT_HIGHLIGHT,
				Placeholder: styles.S_TEXT_DISABLED,
			},
			Blurred: textinput.StyleState{
				Text:        styles.S_TEXT_DISABLED,
				Placeholder: styles.S_TEXT_DISABLED,
			},
			Cursor: styles.TI_CURSOR,
		})

		switch i {
		case 0:
			inp.Placeholder = "Username"
			inp.Validate = func(s string) error {
				if s == "" {
					return fmt.Errorf("too short mate")
				}
				return nil
			}
			inp.Focus()
		case 1:
			inp.Placeholder = "Password"
			inp.EchoMode = textinput.EchoPassword
			inp.Validate = func(s string) error {
				if len(s) < 10 {
					return fmt.Errorf("too short mate")
				}

				return nil
			}
			inp.Blur()
		}
	}

	return &Model{
		focusedField: 0,
		inpName:      inpName,
		inpPass:      inpPass,
	}
}

func (s Model) Init() tea.Cmd {
	return textinput.Blink
}

func (s Model) View() (string, *tea.Cursor) {
	name := s.inpName.View()
	pass := s.inpPass.View()
	btnStyle := styles.STYLE_BTN
	fieldStyle := styles.STYLE_FIELD
	if s.focusedField == 2 {
		if s.inpName.Err != nil || s.inpPass.Err != nil {
			btnStyle = btnStyle.Inherit(STYLE_MOD_DISABLED)
		} else if s.state != 1 {
			btnStyle = styles.STYLE_BTN_SELECTED
		}
	}

	var c *tea.Cursor

	if s.state == 0 && s.focusedField != 2 {
		if s.focusedField == 0 {
			c = s.inpName.Cursor()
		} else {
			c = s.inpPass.Cursor()
		}

		if c != nil {
			c.X += 2
			c.Y += 1
			if s.focusedField == 1 {
				c.Y += 4
			}
		}
	}

	switch s.state {
	case 1:
		btnStyle = btnStyle.Inherit(STYLE_MOD_DISABLED)
		fieldStyle = fieldStyle.Inherit(STYLE_MOD_DISABLED)
	case 2:
		btnStyle = btnStyle.Inherit(STYLE_MOD_WRONG)
		fieldStyle = fieldStyle.Inherit(STYLE_MOD_WRONG)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		fieldStyle.Render(name),
		"",
		fieldStyle.Render(pass),
		"",
		btnStyle.Render("Login"),
	), c
}

// [username, password]
type LoginEntered [2]string

func overflow(min int, v *int, max int) {
	if *v < min {
		*v = max
	} else if *v > max {
		*v = min
	}
}

func (s *Model) changeField(newField int) tea.Cmd {
	textFields := []*textinput.Model{&s.inpName, &s.inpPass}
	if s.focusedField != 2 {
		textFields[s.focusedField].Blur()
	}

	overflow(0, &newField, 2)
	s.focusedField = newField

	if newField != 2 {
		return textFields[newField].Focus()
	}

	return nil
}

func (s *Model) Update(msg tea.Msg) (utils.Screen, tea.Cmd) {
	batcher := []tea.Cmd{}

	switch m := msg.(type) {
	case tea.KeyPressMsg:
		if s.state != 1 {
			switch m.String() {
			case "tab", "down":
				batcher = append(batcher, s.changeField(s.focusedField+1))
			case "shift+tab", "up":
				batcher = append(batcher, s.changeField(s.focusedField-1))
			case "enter":
				if s.focusedField == 2 && s.inpName.Err == nil && s.inpPass.Err == nil {
					s.state = 1

					return s, func() tea.Msg {
						return LoginEntered([2]string{s.inpName.Value(), s.inpPass.Value()})
					}
				} else {
					batcher = append(batcher, s.changeField(s.focusedField+1))
				}
			}
		}
	case clearWrongPass:
		s.state = 0
	}

	var tmpCmd tea.Cmd

	s.inpName, tmpCmd = s.inpName.Update(msg)
	batcher = append(batcher, tmpCmd)

	s.inpPass, tmpCmd = s.inpPass.Update(msg)
	batcher = append(batcher, tmpCmd)

	return s, tea.Batch(batcher...)
}

type clearWrongPass bool

func (s *Model) WrongPassword() tea.Cmd {
	s.state = 2
	s.inpName.SetValue("")
	s.inpPass.SetValue("")
	s.changeField(0)

	return func() tea.Msg {
		<-time.NewTimer(750 * time.Millisecond).C
		return clearWrongPass(true)
	}
}
