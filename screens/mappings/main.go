package mappings

import (
	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/bubbles/textinput"
)

type mappingImpl struct {
	api               *api.APIClient
	fetchedCategories []*api.Category
	categoryField     *textinput.Model
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

func (m *mappingImpl) Init() listeditor.AbstractionSetup {
	v, err := m.api.CategoriesFetch()
	if err != nil {
		panic(err)
	}

	return listeditor.AbstractionSetup{v}
}

func (m *mappingImpl) Setup(msg listeditor.AbstractionSetup) {
	data := msg.V.([]*api.Category)
	m.fetchedCategories = data
	sl := make([]string, len(data))
	for i, v := range data {
		sl[i] = v.Name
	}

	m.categoryField.SetSuggestions(sl)
}

func New(c *api.APIClient, w, h int) *listeditor.Model[mappingProxy, *mappingProxy] {
	m := listeditor.New[mappingProxy](
		"New Mapping", mappingDelegate{}, w, h,
	)
	m.Abstraction = &mappingImpl{
		api: c,
	}

	return m
}
