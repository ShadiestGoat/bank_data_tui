package categories

import (
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/bank_data_tui/utils/listeditor"
)

type categoryDelegate struct{}

func (c categoryDelegate) Spacing() int { return 1 }
func (c categoryDelegate) Height() int  { return 1 }

func (c categoryDelegate) Render(w io.Writer, m list.Model, i int, v list.Item) {
	style := lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	if m.GlobalIndex() == i {
		style = style.Underline(true)
	}

	txt, ok := v.(listeditor.NewItem)
	if ok {
		w.Write(
			[]byte(" " + style.Render(string(txt))),
		)

		return
	}

	cat := v.(*categoryProxy)
	if err := verifyColor(cat.Color); err == nil {
		style = style.Foreground(lipgloss.Color("#" + cat.Color))
	}

	w.Write(
		[]byte(" " + style.Render(utils.Overflow("["+cat.Icon+"] "+cat.Name, listeditor.WIDTH_LIST - 1))),
	)
}

func (c categoryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
