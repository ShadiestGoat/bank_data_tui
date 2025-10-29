package categories

import (
	"slices"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/bank_data_tui/utils/repo"
	tea "github.com/charmbracelet/bubbletea"
)

type categoryImpl struct {
	api *api.APIClient

	cache *repo.Cache
}

func (m *categoryImpl) InitialFetch() ([]*categoryProxy, error) {
	c, err := m.cache.EasyCategories(m.api)
	if err != nil {
		return nil, err
	}

	arr := make([]*categoryProxy, len(c))
	for i, v := range c {
		arr[i] = (*categoryProxy)(v)
	}

	return arr, nil
}

func (m *categoryImpl) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case editor.ItemDel:
		i := slices.IndexFunc(m.cache.Categories, func(c *api.Category) bool {
			return c.ID == string(msg)
		})
		if i != -1 {
			m.cache.Categories = slices.Delete(m.cache.Categories, i, i + 1)
		}
	case listeditor.ItemNew:
		m.cache.Categories = append(m.cache.Categories, (*api.Category)(msg.Value.(*categoryProxy)))
	case listeditor.ItemUpdate:
		cat := (*api.Category)(msg.Value.(*categoryProxy))
		i := slices.IndexFunc(m.cache.Categories, func(c *api.Category) bool {
			return c.ID == string(cat.ID)
		})
		if i == -1 {
			m.cache.Categories = append(m.cache.Categories, cat)
		} else {
			m.cache.Categories[i] = cat
		}
	}
}

func New(c *api.APIClient, cache *repo.Cache, w, h int) *listeditor.Model[categoryProxy, *categoryProxy] {
	m := listeditor.New[categoryProxy](
		"New Category", categoryDelegate{}, w, h,
	)
	m.Abstraction = &categoryImpl{
		api:   c,
		cache: cache,
	}

	return m
}
