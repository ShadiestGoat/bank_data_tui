package filepicker

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) handleKeyCommon(key tea.KeyPressMsg) (tea.Cmd, bool) {
	switch k := key.String(); k {
	case "alt+up":
		if m.acceptedPath == "/" {
			return nil, true
		}

		i := strings.LastIndexByte(m.acceptedPath[:len(m.acceptedPath)-1], '/')
		np := m.acceptedPath[:i]
		if np == "" {
			np = "/"
		}
		cmd := m.forceUserSel(np)
		return cmd, true
	case "shift+up", "shift+down", "shift+tab":
		return m.toggleInput(), true
	}

	return nil, false
}

func (m *Model) handleKeyTextinput(key tea.KeyPressMsg) (tea.Cmd, bool) {
	switch k := key.String(); k {
	case "up", "down":
		if m.sugIndex != -1 {
			sugs := m.suggestions()
			m.sugIndex = cycle(m.sugIndex, len(sugs), k == "down")
		}
		return nil, true
	case "tab", "enter":
		if k == "enter" && m.inpIsFile {
			return func() tea.Msg {
				return FileSelected{m.currentCleanInput()}
			}, true
		}
		if m.sugIndex != -1 {
			sugs := m.suggestions()

			return m.forceUserSel(m.acceptedPath + "/" + sugs[m.sugIndex]), true
		}

		return m.toggleInput(), true
	}

	return nil, false
}

func (m *Model) handleKeyPicker(key tea.KeyPressMsg) (tea.Cmd, bool) {
	switch k := key.String(); k {
	case "up", "down":
		if m.fileIndex != -1 {
			m.fileIndex = cycle(
				m.fileIndex,
				len(m.visibleEntries(m.dirs))+len(m.visibleEntries(m.files)),
				k == "down",
			)
			m.adjustVP()
		}
		return nil, true
	case "tab", "shift+tab":
		return m.toggleInput(), true
	case "enter":
		dirs, files := m.visibleEntries(m.dirs), m.visibleEntries(m.files)
		if m.fileIndex == -1 {
			return nil, true
		}
		if m.fileIndex < len(dirs) {
			return m.forceUserSel(m.acceptedPath + "/" + dirs[m.fileIndex].Name()), true
		}
		return func() tea.Msg {
			return FileSelected{m.acceptedPath + "/" + files[m.fileIndex-len(dirs)].Name()}
		}, true
	case "ctrl+down":
		m.vpOffset++
		m.adjustVP()
	case "ctrl+up":
		m.vpOffset--
		m.adjustVP()
	}

	return nil, false
}
