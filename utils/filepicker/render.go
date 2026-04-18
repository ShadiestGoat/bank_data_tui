package filepicker

import (
	"image/color"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
)

func (m Model) View() (string, *tea.Cursor) {
	if m.w == 0 || m.h == 0 || m.acceptedPath == "" {
		return "", nil
	}

	resp := &strings.Builder{}
	resp.Grow(m.h * (m.w + 1))

	dirs := m.visibleEntries(m.dirs)

	m.viewDirInput(resp)
	resp.WriteRune('\n')
	resp.WriteRune('\n')
	m.viewFiles(resp, dirs)

	var cur *tea.Cursor
	if m.textField.Focused() {
		cur = m.textField.Cursor()
		cur.X += 2
		cur.Y += 1
	}

	return strings.TrimSuffix(resp.String(), "\n"), cur
}

func (m Model) viewDirInput(sb *strings.Builder) {
	text := m.textField.View()
	var textboxColor color.Color = styles.COLOR_DISABLED
	if m.textField.Focused() {
		if m.textField.Value() != "" && m.textField.Position() == lipgloss.Width(m.textField.Value()) {
			text = text[:len(text)-1]
		}
	
		if m.sugIndex != -1 {
			_, f := path.Split(m.textField.Value())
			sug := strings.TrimPrefix(m.suggestions()[m.sugIndex], f)
			text += lipgloss.NewStyle().Faint(true).Foreground(styles.COLOR_DISABLED).Render(sug)
		}

		if m.err != nil {
			textboxColor = styles.COLOR_WRONG
		} else if m.inpIsFile {
			textboxColor = styles.COLOR_SECONDARY
		} else {
			textboxColor = styles.COLOR_MAIN
		}
	}

	sb.WriteString(styles.STYLE_FIELD.Width(m.w).BorderForeground(textboxColor).Render(text))
}

func (m Model) viewFiles(resp *strings.Builder, dirs []os.DirEntry) {
	styleSel := lipgloss.NewStyle().Bold(true).Foreground(styles.COLOR_MAIN)
	styleAllowed := lipgloss.NewStyle().Foreground(styles.COLOR_SECONDARY)
	styleDis := lipgloss.NewStyle().Faint(true).Foreground(styles.COLOR_DISABLED)

	focused := !m.textField.Focused()

	maxExtSize := 0
	maxFileSize := 0

	files := m.visibleEntries(m.files)
	leftH := m.h - 4

	// TODO: Does this need to be memoized? I doubt it matters too much but maybe?
	for _, v := range files {
		ext := filepath.Ext(v.Name())
		if ext != "" {
			if w := lipgloss.Width(ext[1:]); w > maxExtSize {
				maxExtSize = min(w, int(float64(m.w)*0.2))
			}
		}
		if w := lipgloss.Width(v.Name()); w > maxFileSize {
			maxFileSize = w
		}
	}

	for _, v := range dirs {
		// +1 bc of the /
		if w := lipgloss.Width(v.Name()) + 1; w > maxFileSize {
			maxFileSize = w
		}
	}

	maxFileSize = min(maxFileSize, m.w-2-maxExtSize-1)

	for i, v := range dirs[min(m.vpOffset, len(dirs)):] {
		if leftH == 0 {
			return
		}
		style := styleAllowed
		if !focused {
			style = styleDis
		}

		if m.fileIndex == i + m.vpOffset {
			if focused {
				style = styleSel
			}
			resp.WriteString(style.Render("> "))
		} else {
			resp.WriteString(style.Render("  "))
		}

		resp.WriteString(style.Render(
			utils.Overflow(
				strings.Repeat(" ", maxExtSize)+" "+v.Name()+"/",
				m.w-2,
			),
		) + "\n")
		leftH--
	}

	if leftH == 0 {
		return
	}

	off := m.vpOffset - len(dirs)

	if len(dirs) != 0 && len(files) != 0 && off <= 0 {
		resp.WriteString(
			styleDis.Render(
				strings.Repeat(" ", 2+maxExtSize+1)+
					strings.Repeat("─", maxFileSize),
			) + "\n",
		)
		leftH--
	}

	off--
	if off < 0 {
		off = 0
	}

	for i, v := range files[min(off, len(files)):] {
		if leftH == 0 {
			return
		}
		ext := filepath.Ext(v.Name())
		ext = strings.TrimPrefix(ext, ".")

		style := styleAllowed
		if !focused {
			style = styleDis
		}

		if m.fileIndex == i+len(dirs)+off {
			if focused {
				style = styleSel
			}

			resp.WriteString(style.Render("> "))
		} else {
			if len(m.highlightedExt) != 0 && !slices.Contains(m.highlightedExt, ext) {
				style = styleDis
			}

			resp.WriteString(style.Render("  "))
		}

		resp.WriteString(style.Render(
			utils.Overflow(ext, maxExtSize)+strings.Repeat(" ", maxExtSize-lipgloss.Width(ext))+" "+
				utils.Overflow(v.Name(), m.w-1-maxExtSize-2),
		) + "\n")
		leftH--
	}
}
