package api

type SavableCategory struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
	Name  string `json:"name"`
}

type Category struct {
	ID string `json:"id"`

	SavableCategory
}

func (c *APIClient) CategoriesFetch() ([]*Category, error) {
	return deArray(easyFetch[[]*Category](c, `GET`, `/categories`, nil))
}

type RespCreated struct {
	ID string `json:"id"`
}

func (c *APIClient) CategoriesCreate(s *SavableCategory) (string, error) {
	resp, err := easyFetch[RespCreated](c, `POST`, `/categories`, s)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *APIClient) CategoriesUpdate(id string, s *SavableCategory) error {
	return easyNilFetch(c, `PUT`, `/categories/` + id, s)
}
