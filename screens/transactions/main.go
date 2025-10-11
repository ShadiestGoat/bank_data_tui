package transactions

import tea "github.com/charmbracelet/bubbletea"

type Model struct {}

func (v Model) Init() tea.Cmd {return nil}
func (v Model) View() string {return ""}
func (v *Model) Update(tea.Msg) (tea.Model, tea.Cmd) {return v, nil}
