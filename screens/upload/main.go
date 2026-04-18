package upload

import (
	"log"
	"os"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/bank_data_tui/utils/filepicker"
)

type Model struct {
	api           *api.APIClient
	filepicker    filepicker.Model
	uploadingPath string
	err           error
	spin          spinner.Model
	w, h          int
}

const INP_PADDING = 5

type uploaded struct {
	err error
}

func New(api *api.APIClient, w, h int) *Model {
	m := &Model{
		api: api,
		w:   w, h: h,
		spin: spinner.New(spinner.WithStyle(styles.S_TEXT_HIGHLIGHT)),
	}

	fp := filepicker.New(w, h, []string{"tsv", "csv"})

	m.filepicker = fp

	return m
}

func (m Model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m Model) View() (string, *tea.Cursor) {
	box := lipgloss.NewStyle().Width(m.w).Height(m.h).Align(lipgloss.Left, lipgloss.Top)

	if m.uploadingPath == "" {
		res, cur := m.filepicker.View()
		return box.Render(res), cur
	}

	var res string
	if m.err != nil {
		res = lipgloss.JoinVertical(
			lipgloss.Center,
			"!! "+styles.S_TEXT_WRONG.Render("Error Uploading")+" !!",
			m.uploadingPath,
			"",
			styles.S_TEXT_WRONG.Render(m.err.Error()),
		)
	} else {
		spin := m.spin.View()
		res = lipgloss.JoinVertical(
			lipgloss.Center,
			spin+" "+styles.S_TEXT_HIGHLIGHT_SECONDARY.Render("Uploading...")+" "+spin,
			"",
			m.uploadingPath,
		)
	}

	return box.AlignHorizontal(lipgloss.Center).Render(res), nil
}

func (m Model) Update(msg tea.Msg) (utils.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.ResizeMessage:
		m.w, m.h = msg.W, msg.H
		m.filepicker.SetSize(msg.W, msg.H)
		return m, nil
	case uploaded:
		var cmd tea.Cmd
		if msg.err != nil {
			m.err = msg.err
			cmd = clearErrCMD
		} else {
			cmd = utils.GoToHome
		}

		return m, cmd
	case clearErr:
		m.uploadingPath = ""
		m.err = nil
		return m, nil
	case filepicker.FileSelected:
		log.Println("Hey hi!!", msg.Path)
		m.uploadingPath = msg.Path

		return m, tea.Batch(func() tea.Msg {
			f, err := os.Open(msg.Path)
			if err != nil {
				return uploaded{err: err}
			}
			defer f.Close()

			err = m.api.UploadTSV(f)
			return uploaded{err: err}
		}, m.spin.Tick)
	}

	if m.uploadingPath == "" {
		fp, cmd := m.filepicker.Update(msg)
		m.filepicker = fp

		return m, cmd
	} else {
		spin, cmd := m.spin.Update(m)
		m.spin = spin
		return m, cmd
	}
}

type clearErr struct {}
func clearErrCMD() tea.Msg {
	<- time.After(5 * time.Second)
	return clearErr{}
}
