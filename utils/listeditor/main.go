package listeditor

import (
	"github.com/bank_data_tui/utils/editor"
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
	STYLE_SPLIT = lipgloss.NewStyle().MarginLeft(2).PaddingLeft(2).BorderLeft(true).BorderStyle(lipgloss.DoubleBorder())
)

type Abstraction[T any] interface {
	NewEditor(w, h int, curItem T) *editor.Model
	InitialFetch() ([]T, error)
}

type Item interface {
	GetID() string
	SetID(v string)
	FilterValue() string
}

type Model[item any, PT interface {
	Item
	*item
}] struct {
	Abstraction[PT]

	list     list.Model
	spin     spinner.Model
	isLoaded bool
	newItem  NewItem

	items   []PT
	curItem PT

	editor *editor.Model

	w, h int
}

type NewItem string

func (ni NewItem) FilterValue() string { return string(ni) }

func New[T any, PT interface {
	Item
	*T
}](newItemText string, delegate list.ItemDelegate, w, h int) *Model[T, PT] {
	m := &Model[T, PT]{
		spin:     spinner.Model{},
		isLoaded: false,
		newItem:  NewItem(newItemText),
		items:    []PT{},
		curItem:  new(T),
		editor:   &editor.Model{},
		w:        w,
		h:        h,
	}

	m.list = list.New([]list.Item{m.newItem}, delegate, WIDTH_LIST, h)
	m.list.KeyMap = listKeyMap
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)

	return m
}

type initialResp[T any] []T
type AbstractionSetup struct {V any}

func (m *Model[T, PT]) Init() tea.Cmd {
	m.resetEditor()

	batcher := []tea.Cmd{
		func() tea.Msg {
			res, err := m.InitialFetch()
			if err != nil {
				panic("Can't do initial fetch: " + err.Error())
			}

			return initialResp[PT](res)
		},
		m.spin.Tick,
		m.editor.Init(),
	}

	if a, ok := m.Abstraction.(interface {Init() AbstractionSetup}); ok {
		batcher = append(batcher, func() tea.Msg {
			return a.Init()
		})
	}

	return tea.Batch(batcher...)
}

func (m Model[T, PT]) View() string {
	if !m.isLoaded {
		return m.spin.View()
	}

	l, e := m.list.View(), m.editor.View()

	res := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(WIDTH_LIST).AlignHorizontal(lipgloss.Left).Render(l),
		STYLE_SPLIT.Height(m.h).Render(e),
	)

	return res
}

func (m *Model[T, PT]) Resize(w, h int) {
	m.w, m.h = w, h

	m.list.SetHeight(h)
	m.editor.SetWidth(w - WIDTH_OFFSET_EDITOR)
}
