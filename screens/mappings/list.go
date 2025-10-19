package mappings

import (
	"io"

	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mappingDelegate struct{}

func (mappingDelegate) Spacing() int { return 1 }
func (mappingDelegate) Height() int  { return 1 }

func (mappingDelegate) Render(w io.Writer, m list.Model, i int, v list.Item) {
	style := lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	if m.GlobalIndex() == i {
		style = style.Underline(true)
	}

	txt, ok := v.(listeditor.NewItem)
	if ok {
		w.Write(
			[]byte(" " + style.Render("| "+string(txt))),
		)

		return
	}

	val := v.(*mappingProxy)
	w.Write([]byte(" " + style.Render(val.Name)))
}

func (mappingDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
