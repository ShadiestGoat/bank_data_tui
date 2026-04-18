package mappings

import (
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/bank_data_tui/utils/listeditor"
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
	w.Write([]byte(" " + style.Render(
		utils.Overflow(val.Name, listeditor.WIDTH_LIST - 1),
	)))
}

func (mappingDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
