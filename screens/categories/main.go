package categories

import (
	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/listeditor"
)

type categoryImpl struct {
	api *api.APIClient
}

func (m *categoryImpl) InitialFetch() ([]*categoryProxy, error) {
	c, err := m.api.CategoriesFetch()
	if err != nil {
		return nil, err
	}
	arr := make([]*categoryProxy, len(c))
	for i, v := range c {
		arr[i] = (*categoryProxy)(v)
	}

	return arr, nil
}

func New(c *api.APIClient, w, h int) *listeditor.Model[categoryProxy, *categoryProxy] {
	m := listeditor.New[categoryProxy](
		"New Category", categoryDelegate{}, w, h,
	)
	m.Abstraction = &categoryImpl{c}

	return m
}
