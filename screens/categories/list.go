package categories

import (
	"io"

	"github.com/bank_data_tui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type categoryDelegate struct {}

func (c categoryDelegate) Spacing() int { return 1 }
func (c categoryDelegate) Height() int { return 1 }

func (c categoryDelegate) Render(w io.Writer, m list.Model, i int, v list.Item) {
	style := lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	if m.GlobalIndex() == i {
		style = style.Underline(true)
	}

	_, ok := v.(*newCategoryItem)
	if ok {
		w.Write(
			[]byte(" " + style.Render("New Category")),
		)

		return
	}

	cat := v.(*categoryItem)
	if err := verifyColor(cat.Color); err == nil {
		style = style.Foreground(lipgloss.Color("#" + cat.Color))
	}

	w.Write(
		[]byte(" " + style.Render("[" + cat.Icon + "] " + cat.Name)),
	)
}

func (c categoryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

var ListKeymap = list.KeyMap{
	CursorUp: key.NewBinding(
		key.WithKeys("alt+up"),
		key.WithHelp("alt ↑", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("alt+down"),
		key.WithHelp("alt ↓", "down"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("alt+pgup"),
		key.WithHelp("alt page up", "prev page"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("alt+pgdown"),
		key.WithHelp("alt page down", "next page"),
	),
	GoToStart: key.NewBinding(
		key.WithKeys("alt+home"),
		key.WithHelp("alt home", "go to start"),
	),
	GoToEnd: key.NewBinding(
		key.WithKeys("alt+end"),
		key.WithHelp("alt end", "go to end"),
	),
	Filter: key.NewBinding(
		key.WithKeys("alt+/"),
		key.WithHelp("alt /", "filter"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("alt esc", "clear filter"),
	),

	// Filtering.
	CancelWhileFiltering: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("alt esc", "cancel"),
	),
	AcceptWhileFiltering: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("alt enter", "apply filter"),
	),
}

func doesKeyMatchList(k tea.KeyMsg, l list.Model) bool {
	if l.FilterState() == list.Filtering {
		return key.Matches(
			k,
			l.KeyMap.CancelWhileFiltering,
			l.KeyMap.AcceptWhileFiltering,
		)
	} else if l.FilterState() == list.FilterApplied {
		if key.Matches(k, l.KeyMap.ClearFilter) {
			return true
		}
	}

	return key.Matches(
		k,
		l.KeyMap.CursorUp,
		l.KeyMap.CursorDown,
		l.KeyMap.PrevPage,
		l.KeyMap.NextPage,
		l.KeyMap.GoToStart,
		l.KeyMap.GoToEnd,
		l.KeyMap.Filter,
	)
}