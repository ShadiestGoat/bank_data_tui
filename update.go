package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/bank_data_tui/screens/categories"
	"github.com/bank_data_tui/screens/login"
	"github.com/bank_data_tui/screens/mappings"
	"github.com/bank_data_tui/screens/transactions"
	"github.com/bank_data_tui/screens/upload"
	"github.com/bank_data_tui/utils"
)

func (m *mainApp) switchToScreen(s Screen) tea.Cmd {
	if m.curFocusedScreen == s {
		return nil
	}

	m.curFocusedScreen = s
	switch s {
	case S_TRANS:
		m.screenImp = transactions.New(m.api, m.cache, m.width, m.height-HEADER_HEIGHT)
	case S_MAPPINGS:
		m.screenImp = mappings.New(m.api, m.cache, m.width, m.height-HEADER_HEIGHT)
	case S_CATEGORIES:
		m.screenImp = categories.New(m.api, m.cache, m.width, m.height-HEADER_HEIGHT)
	case S_UPLOAD:
		m.screenImp = upload.New(m.api, m.width, m.height-HEADER_HEIGHT)
	}

	return m.screenImp.Init()
}

func (m *mainApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	batcher := []tea.Cmd{}

	passToChildren := false

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "alt+tab":
			s := m.curFocusedScreen + 1
			if s > S_UPLOAD {
				s = S_TRANS
			}
			batcher = append(batcher, m.switchToScreen(s))
		case "alt+shift+tab":
			s := m.curFocusedScreen - 1
			if s == S_LOGIN {
				s = S_UPLOAD
			}
			batcher = append(batcher, m.switchToScreen(s))
		case "alt+t":
			batcher = append(batcher, m.switchToScreen(S_TRANS))
		case "alt+m":
			batcher = append(batcher, m.switchToScreen(S_MAPPINGS))
		case "alt+c":
			batcher = append(batcher, m.switchToScreen(S_CATEGORIES))
		case "alt+u", "alt+n":
			batcher = append(batcher, m.switchToScreen(S_UPLOAD))
		default:
			passToChildren = true
		}
	case tea.WindowSizeMsg:
		log.Println("RESIZE", msg.Width)
		m.height = msg.Height
		m.width = msg.Width

		m.screenImp, cmd = m.screenImp.Update(utils.ResizeMessage{
			W: m.width,
			H: m.height - HEADER_HEIGHT,
		})
		batcher = append(batcher, cmd)
	case login.LoginEntered:
		screen, ok := m.screenImp.(*login.Model)
		if !ok {
			panic("Somehow on the wrong model?")
		}

		err := m.api.Login(msg)
		if err != nil {
			batcher = append(batcher, screen.WrongPassword())
		} else {
			batcher = append(batcher, m.switchToScreen(S_TRANS))
		}
	case utils.MsgGoToHome:
		batcher = append(batcher, m.switchToScreen(S_TRANS))
	default:
		passToChildren = true
	}

	if passToChildren {
		m.screenImp, cmd = m.screenImp.Update(msg)
		batcher = append(batcher, cmd)
	}

	return m, tea.Batch(batcher...)
}
