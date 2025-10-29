package transactions

import (
	"log"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/repo"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	w, h        int
	selected    int
	viewportOff int
	items       []*api.Transaction
	api         *api.APIClient
	cache       *repo.Cache
}

func New(api *api.APIClient, cache *repo.Cache, w, h int) *Model {
	return &Model{
		w:     w,
		h:     h,
		api:   api,
		cache: cache,
	}
}

type newPageData struct{ data []*api.Transaction }

func (m Model) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		d, err := m.api.TransactionsFetch(api.TOR_AUTH, 1, false)
		if err != nil {
			panic(err)
		}

		return newPageData{d.Data}
	}, func() tea.Msg {
		_, err := m.cache.EasyCategories(m.api)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.selected = m.overflow(m.selected + 1)
		case "up":
			m.selected = m.overflow(m.selected - 1)
		case "end":
			m.selected = len(m.items) - 1
		case "start":
			m.selected = 0
		case "alt+down":
			m.viewportOff = m.viewportOff + 1
			visibleItems := len(m.items) - m.viewportOff
			if m.h - visibleItems > 8 {
				m.viewportOff--
			}
			// if len(m.items) - m.viewportOff < 7 {
			// 	m.viewportOff--
			// }
		case "alt+up":
			m.viewportOff = m.viewportOff - 1
			if m.viewportOff < 0 {
				m.viewportOff = 0
			}
		}

		if msg.Alt {
			m.forceSelIntoViewport()
		} else {
			m.forceViewportIntoSel()
		}
	case newPageData:
		m.items = msg.data
	}

	return m, nil
}

func (m *Model) overflow(v int) int {
	if v < 0 {
		return len(m.items) - 1
	} else if v > len(m.items) - 1 {
		return 0
	}

	return v
}

func (m *Model) forceViewportIntoSel() {
	if len(m.items) <= m.h {
		m.viewportOff = 0
		return
	}

	log.Println(m.selected, m.viewportOff, m.h)
	if m.selected < m.viewportOff {
		m.viewportOff = m.selected
	} else if m.selected > m.viewportOff + m.h - 2 {
		m.viewportOff = m.selected - m.h + 2
	}
}

func (m *Model) forceSelIntoViewport() {
	if len(m.items) <= m.h {
		m.viewportOff = 0
		return
	}

	if m.selected < m.viewportOff {
		m.selected = m.viewportOff
	} else if m.selected > m.viewportOff + m.h - 2 {
		m.selected = m.viewportOff + m.h - 2
	}
}

func (m *Model) Resize(w, h int) {
	m.w, m.h = w, h
	m.forceViewportIntoSel()
}
