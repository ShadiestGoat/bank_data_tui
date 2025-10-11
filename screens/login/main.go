package login

import (
	"fmt"
	"time"

	"github.com/bank_data_tui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	STYLE_MOD_OK       = lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	STYLE_MOD_DISABLED = lipgloss.NewStyle().BorderForeground(lipgloss.ANSIColor(8)).Background(lipgloss.ANSIColor(8))
	STYLE_MOD_WRONG    = lipgloss.NewStyle().BorderForeground(lipgloss.ANSIColor(9)).Background(lipgloss.ANSIColor(9))
)

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
		inp.Width = 15
		inp.Prompt = ""
		inp.TextStyle = lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)

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

func (s Model) View() string {
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

	switch s.state {
	case 0:
		btnStyle = btnStyle.Inherit(STYLE_MOD_OK)
		fieldStyle = fieldStyle.Inherit(STYLE_MOD_OK)
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
	)
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

func (s *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	batcher := []tea.Cmd{}

	switch m := msg.(type) {
	case tea.KeyMsg:
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
