package upload

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/styles"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	api      *api.APIClient
	dirInput textinput.Model
	picker   filepicker.Model
	spin     spinner.Model
	path     string

	home string
	cwd  string

	w, h int
}

const INP_PADDING = 5

var fileErr = errors.New("is file :?")

func New(api *api.APIClient, w, h int) *Model {
	fp := filepicker.New()
	fp.SetHeight(h - 5)
	fp.AllowedTypes = []string{"tsv", "csv"}
	fp.DirAllowed = false
	fp.ShowPermissions = false

	inp := textinput.New()
	inp.Width = w - 5
	inp.Prompt = ""
	inp.Focus()
	inp.ShowSuggestions = true

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	m := &Model{
		api:    api,
		picker: fp,
		spin:   spinner.New(spinner.WithSpinner(spinner.Ellipsis)),
		w:      w, h: h,
		home: home, cwd: cwd,
	}

	inp.Validate = func(s string) error {
		s = m.cleanPath(s)

		info, err := os.Stat(s)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fileErr
		}

		return nil
	}

	fp.CurrentDirectory = m.cleanPath(".")
	inpDir, ok := strings.CutPrefix(fp.CurrentDirectory, m.home)
	if ok {
		inpDir = "~" + inpDir
	}
	inp.SetValue(inpDir)

	m.dirInput = inp

	return m
}

func (m Model) cleanPath(s string) string {
	suffix := ""
	if strings.HasSuffix(s, "/") {
		suffix = "/"
	}
	if s == "" {
		return "/"
	}

	if s[0] != '/' && s != "~" && !strings.HasPrefix(s, "~/") {
		s = m.cwd + "/" + s
	}

	s = path.Clean(s)
	if s == "~" {
		return m.home
	}
	if s == "/" {
		return s
	}

	s, ok := strings.CutPrefix(s, "~/")
	if ok {
		s = m.home + "/" + s
	}

	return s + suffix
}

func (m *Model) Resize(w, h int) {
	m.w, m.h = w, h
	m.picker.SetHeight(h - 5)
	m.dirInput.Width = w - 5
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.dirInput.Focus(),
		m.picker.Init(),
	)
}

func (m Model) View() string {
	box := lipgloss.NewStyle().Width(m.w).Height(m.h).Align(lipgloss.Left, lipgloss.Right)

	if m.path != "" {
		return box.Render(m.spin.View() + " Uploading transactions")
	}

	inp := lipgloss.NewStyle().MaxWidth(m.dirInput.Width).Render(m.dirInput.View())
	fieldStyle := styles.STYLE_FIELD
	if m.dirInput.Err != nil {
		if m.dirInput.Focused() {
			fieldStyle = fieldStyle.BorderForeground(styles.COLOR_WRONG)
		} else {
			fieldStyle = fieldStyle.BorderForeground(styles.COLOR_DISABLED)
		}
	} else if m.dirInput.Focused() {
		fieldStyle = fieldStyle.BorderForeground(styles.COLOR_MAIN)
	}

	return box.Render(
		fieldStyle.Render(inp) + "\n\n" +
			m.picker.View(),
	)
}

func (m *Model) setSuggestions() {
	if m.dirInput.Err != nil {
		return
	}

	p := m.cleanPath(m.dirInput.Value())
	if strings.HasSuffix(m.dirInput.Value(), "/") {
		p = strings.TrimSuffix(p, "/")
	} else {
		p = path.Dir(p)
	}

	suggBase := p
	if strings.HasPrefix(m.dirInput.Value(), "~/") {
		suggBase = "~" + strings.TrimPrefix(p, m.home)
	}

	dirEntry, err := os.ReadDir(p)
	if err == nil {
		suggs := []string{}
		for _, e := range dirEntry {
			if !strings.HasPrefix(e.Name(), ".") && e.IsDir() {
				suggs = append(suggs, suggBase+"/"+e.Name())
			}
		}

		m.dirInput.SetSuggestions(suggs)
	}
}

type FileUploadComplete bool

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	batcher := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "down", "up":
			if m.dirInput.Focused() {
				sug := m.dirInput.CurrentSuggestion()
				if msg.String() == "tab" && sug != "" && sug != m.dirInput.Value() {
					m.dirInput.SetValue(m.dirInput.CurrentSuggestion())
				} else {
					m.dirInput.Blur()
				}
			} else {
				batcher = append(batcher, m.dirInput.Focus())
			}
		case "enter":
			if m.dirInput.Focused() && m.dirInput.Err == fileErr {
				batcher = append(batcher, m.selectPath(m.dirInput.Value()))
			}
		}
	}

	if m.path != "" {
		m.spin, cmd = m.spin.Update(msg)
		batcher = append(batcher, cmd)
	} else {
		if m.dirInput.Focused() {
			m.dirInput, cmd = m.dirInput.Update(msg)
			if m.dirInput.Err == nil && m.dirInput.Value() != "" {
				np := m.cleanPath(m.dirInput.Value())
				if np != m.picker.CurrentDirectory {
					m.picker.CurrentDirectory = np
					batcher = append(batcher, m.picker.Init())
				}
			}

			batcher = append(batcher, cmd)
			m.setSuggestions()
		}

		_, isKey := msg.(tea.KeyMsg)
		if !isKey || !m.dirInput.Focused() {
			oldPath := m.cleanPath(m.picker.CurrentDirectory)
			m.picker, cmd = m.picker.Update(msg)
			batcher = append(batcher, cmd)

			if oldPath != m.cleanPath(m.picker.CurrentDirectory) {
				m.dirInput.SetValue(m.cleanPath(m.picker.CurrentDirectory))
			}
		}
	}

	if m.path == "" && !m.dirInput.Focused() {
		if didSelect, path := m.picker.DidSelectFile(msg); didSelect {
			batcher = append(batcher, m.selectPath(path))
		}
	}

	return m, tea.Batch(batcher...)
}

func (m *Model) selectPath(p string) tea.Cmd {
	m.path = p
	return tea.Batch(
		m.spin.Tick,
		func() tea.Msg {
			f, err := os.Open(m.cleanPath(p))
			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			err = m.api.UploadTSV(f)
			if err != nil {
				log.Fatal(err)
			}

			return FileUploadComplete(true)
		},
	)
}