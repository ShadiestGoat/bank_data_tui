// My own impl of the file picker, with a text bar on top and without the weirdness
package filepicker

import (
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/bank_data_tui/styles"
)

type Model struct {
	files     []os.DirEntry
	dirs      []os.DirEntry
	fileIndex int
	sugIndex  int
	inpIsFile bool
	err       error
	textField textinput.Model
	home      string
	// Accepted clean path, not ending in a / unless its == /
	acceptedPath   string
	w, h           int
	highlightedExt []string
	vpOffset       int
}

type readDirMsg struct {
	inp         string
	files, dirs []os.DirEntry
	err         error
	isFile      bool
	force       bool
}

type FileSelected struct {
	Path string
}

func readDir(dir string, force bool) *readDirMsg {
	stat, err := os.Stat(dir)
	if err != nil {
		return &readDirMsg{inp: dir, err: err}
	}
	if !stat.IsDir() {
		return &readDirMsg{inp: dir, isFile: true}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return &readDirMsg{inp: dir, err: err}
	}

	files := make([]os.DirEntry, 0, len(entries))
	dirs := make([]os.DirEntry, 0, len(entries))

	for _, v := range entries {
		if v.IsDir() {
			dirs = append(dirs, v)
		} else {
			files = append(files, v)
		}
	}

	return &readDirMsg{
		inp:   dir,
		files: files,
		dirs:  dirs,
		err:   err,
		force: force,
	}
}

func readDirCMD(dir string, force bool) tea.Cmd {
	return func() tea.Msg {
		return readDir(dir, force)
	}
}

func (m Model) cleanPath(inp string) string {
	if inp == "" || inp == "/" {
		return inp
	}
	hasSlash := inp[len(inp)-1] == '/'

	if m.home != "" {
		if inp == "~" {
			return m.home
		}

		if strings.HasPrefix(inp, "~/") {
			inp = m.home + "/" + inp[2:]
		}
	}
	if hasSlash {
		inp = inp[:len(inp)-1]
	}

	inp = os.ExpandEnv(inp)
	inp = path.Clean(inp)

	return inp
}

func (m Model) presentablePath(p string) string {
	if aft, ok := strings.CutPrefix(p, m.home); ok {
		return "~" + aft
	}

	return p
}

func (m Model) visibleEntries(f []os.DirEntry) []os.DirEntry {
	// In the future I might make this configurable so
	out := make([]os.DirEntry, 0, len(f))
	for _, v := range f {
		if v.Name()[0] != '.' {
			out = append(out, v)
		}
	}

	return out
}

func (m *Model) SetSize(w, h int) {
	m.w, m.h = w, h
	m.textField.CharLimit = int(float64(w-4) * 0.75)
}

func (m Model) currentCleanInput() string {
	return m.cleanPath(m.textField.Value())
}

func (m *Model) updateUsedDir(msg *readDirMsg) {
	m.inpIsFile = false
	m.dirs = msg.dirs
	m.files = msg.files
	m.acceptedPath = msg.inp
	m.fileIndex = 0
	m.vpOffset = 0
	if len(m.suggestions()) != 0 {
		m.sugIndex = 0
	} else {
		m.sugIndex = -1
	}
}

func (m *Model) adjustVP() {
	if m.fileIndex == -1 {
		m.vpOffset = 0
		return
	}

	dirCount := len(m.visibleEntries(m.dirs))
	fCount := len(m.visibleEntries(m.files))
	fileRow := m.fileIndex
	if m.fileIndex >= dirCount {
		fileRow++
	}

	// else if m.vpOffset + pickerHeight > fileRow {
	// 	m.vpOffset = fileRow
	// } else if fc := len(m.visibleEntries(m.files)); m.vpOffset + pickerHeight > fc + dirCount {
	// 	m.vpOffset = fc + dirCount - pickerHeight
	// }

	pickerHeight := m.h - 4
	if fCount+dirCount+1 < pickerHeight {
		m.vpOffset = 0
		return
	}

	highestRow := m.vpOffset + pickerHeight - 1

	log.Println("ADJUST PRE", m.vpOffset, fileRow, highestRow)

	if m.vpOffset > fileRow {
		m.vpOffset = fileRow
	} else if m.vpOffset < 0 {
		m.vpOffset = 0
	} else if highestRow < fileRow {
		m.vpOffset += fileRow - highestRow
	} else if highestRow > fCount+dirCount {
		m.vpOffset = fCount + dirCount - pickerHeight + 1
	}

	log.Println("ADJUST POST", m.vpOffset)

}

func cycle(cur int, len int, inc bool) int {
	if inc {
		v := cur + 1
		if v >= len {
			return 0
		}
		return v
	} else {
		if cur == 0 {
			return len - 1
		}

		return cur - 1
	}
}

func (m *Model) toggleInput() tea.Cmd {
	if m.textField.Focused() {
		m.textField.Blur()
		return nil
	}
	return m.textField.Focus()
}

func (m *Model) forceUserSel(dir string) tea.Cmd {
	log.Println("MEOW", dir, path.Clean(dir))
	dir = path.Clean(dir)
	p := m.presentablePath(dir)
	if !strings.HasSuffix(p, "/") {
		p += "/"
	}

	m.textField.SetValue(p)
	m.textField.CursorEnd()
	return readDirCMD(dir, true)
}

func (m *Model) suggestions() []string {
	inp := m.textField.Value()
	if inp == "" {
		return nil
	}
	_, f := path.Split(m.textField.Value())

	possibleSugs := m.visibleEntries(m.dirs)

	sugs := make([]string, 0, len(possibleSugs))
	for _, v := range possibleSugs {
		if strings.HasPrefix(v.Name(), f) {
			sugs = append(sugs, v.Name())
		}
	}

	return sugs
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *readDirMsg:
		if msg.force || m.currentCleanInput() == msg.inp {
			if msg.err != nil {
				m.err = msg.err
			} else {
				m.err = nil
				if m.acceptedPath == "" {
					m.textField.SetValue(m.presentablePath(msg.inp) + "/")
				}

				if msg.isFile {
					m.inpIsFile = true
				} else if m.acceptedPath != msg.inp {
					m.updateUsedDir(msg)
				}
			}
		}
	case tea.KeyPressMsg:
		cmd, ok := m.handleKeyCommon(msg)
		if ok {
			return m, cmd
		}
		if m.textField.Focused() {
			cmd, ok := m.handleKeyTextinput(msg)
			if ok {
				return m, cmd
			}
		} else {
			cmd, ok := m.handleKeyPicker(msg)
			if ok {
				return m, cmd
			}
		}
	}

	if m.textField.Focused() {
		commands := []tea.Cmd{}

		last := m.textField.Value()
		var lastSug string
		if m.sugIndex != -1 {
			lastSug = m.suggestions()[m.sugIndex]
		}

		tf, cmd := m.textField.Update(msg)
		m.textField = tf
		commands = append(commands, cmd)
		if m.acceptedPath != "" && m.textField.Value() == "" {
			m.sugIndex = -1
			commands = append(commands, m.forceUserSel("/"))
			return m, tea.Batch(commands...)
		}

		if cur := m.textField.Value(); last != cur {
			slashed := strings.HasSuffix(m.textField.Value(), "/")

			nd := m.currentCleanInput()
			if !slashed {
				nd = path.Dir(nd)
			}

			if m.acceptedPath != nd && strings.HasPrefix(m.acceptedPath, nd) {
				// this means we moved into a dir above
				m.sugIndex = -1
				commands = append(commands, readDirCMD(nd, true))
			} else {
				sugs := m.suggestions()
				ni := slices.Index(sugs, lastSug)
				if len(sugs) == 0 {
					m.sugIndex = -1
				} else if ni == -1 {
					m.sugIndex = 0
				} else {
					m.sugIndex = ni
				}

				commands = append(commands, readDirCMD(m.cleanPath(cur), false))
			}
		}

		return m, tea.Batch(commands...)
	}

	return m, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		dir, err := os.Getwd()
		if err != nil {
			dir = m.home
		}
		if dir == "" {
			dir = "/"
		}

		return readDir(dir, true)
	}, textinput.Blink)
}

func New(w, h int, extensions []string) Model {
	home, _ := os.UserHomeDir() // err here doesn't matter
	ti := textinput.New()
	ti.Prompt = ""
	ti.SetVirtualCursor(false)
	ti.SetStyles(textinput.Styles{
		Focused: textinput.StyleState{},
		Blurred: textinput.StyleState{
			Text: styles.S_TEXT_DISABLED,
		},
		Cursor: styles.TI_CURSOR,
	})
	ti.Focus()

	m := Model{
		sugIndex:       -1,
		textField:      ti,
		home:           home,
		highlightedExt: extensions,
	}
	m.SetSize(w, h)

	return m
}
