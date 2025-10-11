package main

import (
	"log"
	"os"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/screens/login"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

type Screen int

const (
	S_LOGIN Screen = iota
	S_TRANS
	S_MAPPINGS
	S_CATEGORIES
	S_UPLOAD
)

type mainApp struct {
	curFocusedScreen Screen
	screenImp        tea.Model

	width  int
	height int

	api *api.APIClient
}

func (m mainApp) Init() tea.Cmd {
	return m.screenImp.Init()
}

func main() {
	f, err := os.Create("logs/log.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	defer f.Close()
	godotenv.Load()

	app := &mainApp{
		curFocusedScreen: S_LOGIN,
		screenImp:        login.NewScreenLogin(),
		api:              &api.APIClient{},
	}
	user, pass := os.Getenv("USERNAME"), os.Getenv("PASSWORD")

	if user != "" && pass != "" {
		err := app.api.Login([2]string{user, pass})
		if err != nil {
			panic(err)
		}
		app.switchToScreen(S_TRANS)
	}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic("Shit fuck shit: " + err.Error())
	}
}
