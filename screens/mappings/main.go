package mappings

import (
	"log"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/bank_data_tui/utils/repo"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type mappingImpl struct {
	cache         *repo.Cache
	api           *api.APIClient
	categoryField *textinput.Model
}

func (m *mappingImpl) InitialFetch() ([]*mappingProxy, error) {
	all, err := m.api.MappingsFetch()
	if err != nil {
		return nil, err
	}
	arr := make([]*mappingProxy, len(all))
	for i, v := range all {
		arr[i] = (*mappingProxy)(v)
	}

	return arr, nil
}

type absSetup struct {}

func (m *mappingImpl) Init() tea.Cmd {
	return func() tea.Msg {
		_, err := m.cache.EasyCategories(m.api)
		if err != nil {
			panic(err)
		}
		log.Println("Murka", len(m.cache.Categories))

		return absSetup{}
	}
}

func (m *mappingImpl) Update(msg tea.Msg) {
	switch msg.(type) {
	case absSetup:
		m.resetSuggestions()
	}
}

func (m *mappingImpl) resetSuggestions() {
	sl := make([]string, len(m.cache.Categories))
	for i, v := range m.cache.Categories {
		sl[i] = v.Name
	}

	m.categoryField.SetSuggestions(sl)
}

func New(c *api.APIClient, cache *repo.Cache, w, h int) *listeditor.Model[mappingProxy, *mappingProxy] {
	m := listeditor.New[mappingProxy](
		"New Mapping", mappingDelegate{}, w, h,
	)
	m.Abstraction = &mappingImpl{
		api:   c,
		cache: cache,
	}

	return m
}
