package categories

import (
	"io"

	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		[]byte(" " + style.Render("["+cat.Icon+"] "+cat.Name)),
	)
}

func (c categoryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
