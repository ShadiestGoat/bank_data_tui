package categories

import (
	"log"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	WIDTH_LIST               = 25
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

type categoryItem struct {
	name string
}

func (c categoryItem) FilterValue() string { return c.name }
func (c categoryItem) Title() string       { return c.name }
func (c categoryItem) Description() string { return "" }

type newCategoryItem struct{}

func (c newCategoryItem) FilterValue() string { return string([]byte{0}) }
func (c newCategoryItem) Title() string       { return "Make a new one" }
func (c newCategoryItem) Description() string { return "Special Option!" }

func New(c *api.APIClient, w, h int) *Model {
	list := list.New([]list.Item{
		newCategoryItem{},
	}, list.NewDefaultDelegate(), WIDTH_LIST, h)

	list.KeyMap.Quit.SetEnabled(false)
	list.SetShowTitle(false)
	list.SetShowHelp(false)

	ni := &api.Category{}

	return &Model{
		c:    c,
		spin: spinner.New(spinner.WithSpinner(spinner.Points)),
		list: list,
		w:    w, h: h,
		curItem: ni,
		editor:  newCategoryEditor(c, w, ni),
	}
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
	switch msg := msg.(type) {
	case resp:
		m.categories = msg
		m.isLoaded = true
	}

	batcher := []tea.Cmd{}
	var cmd tea.Cmd

	if !m.isLoaded {
		m.spin, cmd = m.spin.Update(msg)
		batcher = append(batcher, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	batcher = append(batcher, cmd)

	m.editor, cmd = m.editor.Update(msg)
	batcher = append(batcher, cmd)

	return m, tea.Batch(batcher...)
}
