package categories

import (
	"log"
	"slices"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	WIDTH_LIST               = 15
	WIDTH_EDITOR_SPLIT_SPACE = 2
	WIDTH_OFFSET_EDITOR      = WIDTH_LIST + 1 + WIDTH_EDITOR_SPLIT_SPACE + WIDTH_EDITOR_SPLIT_SPACE // border + margin + padding
)

var (
	STYLE_EDITOR = lipgloss.NewStyle().MarginLeft(2).PaddingLeft(2).BorderLeft(true).BorderStyle(lipgloss.DoubleBorder())
)

type Model struct {
	list       list.Model
	spin       spinner.Model
	isLoaded   bool
	c          *api.APIClient
	categories []*api.Category

	editor  *editor.Model
	curItem *api.Category

	w, h int
}

type categoryItem struct{ *api.SavableCategory }

func (c categoryItem) FilterValue() string { return c.Name + c.Icon }

type newCategoryItem struct{}

func (c newCategoryItem) FilterValue() string { return string([]byte{0}) }

func New(c *api.APIClient, w, h int) *Model {
	listM := list.New([]list.Item{
		&newCategoryItem{},
	}, &categoryDelegate{}, WIDTH_LIST, h)

	listM.KeyMap = list.KeyMap{
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
	listM.SetShowTitle(false)
	listM.SetShowHelp(false)

	ni := &api.Category{}

	m := &Model{
		c:    c,
		spin: spinner.New(spinner.WithSpinner(spinner.Points)),
		list: listM,
		w:    w, h: h,
		curItem: ni,
	}

	m.resetCategoryEditor()

	return m
}

type resp []*api.Category

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			res, err := m.c.CategoriesFetch()
			log.Println("FetchCats:", err, res)

			return resp(res)
		},
		m.spin.Tick,
		m.editor.Init(),
	)
}

func (m Model) View() string {
	if !m.isLoaded {
		return m.spin.View()
	}

	l, e := m.list.View(), m.editor.View()

	res := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(WIDTH_LIST).AlignHorizontal(lipgloss.Left).Render(l),
		STYLE_EDITOR.Height(m.h).Render(e),
	)

	return res
}

func (m *Model) Resize(w, h int) {
	m.w, m.h = w, h

	m.list.SetHeight(h)
	m.editor.SetWidth(w - WIDTH_OFFSET_EDITOR)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	batcher := []tea.Cmd{}
	var cmd tea.Cmd
	bubble := true

	switch msg := msg.(type) {
	case resp:
		m.categories = msg
		cmd = m.list.SetItems(m.categoryItems())
		batcher = append(batcher, cmd)
		m.isLoaded = true
	case editor.ItemNew:
		m.curItem.ID = string(msg)
		m.categories = append(m.categories, m.curItem)
		batcher = append(batcher, m.list.SetItems(m.categoryItems()))
	case editor.ItemDel:
		i := slices.IndexFunc(m.categories, func(c *api.Category) bool { return c.ID == string(msg) })
		if i != -1 {
			m.categories = slices.Delete(m.categories, i, i +1)
		}
		batcher = append(batcher, m.list.SetItems(m.categoryItems()))
	case tea.KeyMsg:
		switch msg.String() {
		case "alt+up":
			bubble = false
			m.list.CursorUp()
		case "alt+down":
			bubble = false
			m.list.CursorDown()
		}
	}

	if !m.isLoaded {
		m.spin, cmd = m.spin.Update(msg)
		batcher = append(batcher, cmd)
	}

	if bubble {
		forList := false
		if msg, ok := msg.(tea.KeyMsg); ok {
			forList = doesKeyMatchList(msg, m.list)
		}

		m.list, cmd = m.list.Update(msg)
		batcher = append(batcher, cmd)

		if !forList && m.list.FilterState() != list.Filtering {
			m.editor, cmd = m.editor.Update(msg)
			batcher = append(batcher, cmd)
		}
	}

	i := m.list.GlobalIndex()
	log.Println("GI", i, len(m.categories))
	if m.isNewCategory(i) {
		if m.curItem.ID != "" {
			m.curItem = &api.Category{}
			cmd = m.resetCategoryEditor()
			batcher = append(batcher, cmd)
		}
	} else if m.categories[i].ID != m.curItem.ID {
		m.curItem = m.categories[i]
		cmd = m.resetCategoryEditor()
		batcher = append(batcher, cmd)
	}

	return m, tea.Batch(batcher...)
}

func (m *Model) isNewCategory(gi int) bool {
	return gi >= len(m.categories)
}

func (m *Model) categoryItems() []list.Item {
	arr := make([]list.Item, len(m.categories)+1)
	arr[len(arr)-1] = &newCategoryItem{}
	for i, v := range m.categories {
		arr[i] = &categoryItem{&v.SavableCategory}
	}

	return arr
}
