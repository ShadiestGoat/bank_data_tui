package repo

import "github.com/bank_data_tui/api"

type Cache struct {
	Categories []*api.Category
}

func (s *Cache) EasyCategories(c *api.APIClient) ([]*api.Category, error) {
	if s.Categories != nil {
		return s.Categories, nil
	}

	v, err := c.CategoriesFetch()
	if err != nil {
		return nil, err
	}

	s.Categories = v
	return v, nil
}
