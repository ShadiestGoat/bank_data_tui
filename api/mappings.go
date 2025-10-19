package api

type Mapping struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`

	InpText string  `json:"inputText,omitempty"`
	InpAmt  *float64 `json:"inputAmount,omitempty"`

	ResName       string `json:"resName,omitempty"`
	ResCategoryID string `json:"resCategoryID,omitempty"`

	Priority int `json:"priority"`
}

func (c *APIClient) MappingsFetch() ([]*Mapping, error) {
	return deArray(easyFetch[[]*Mapping](c, `GET`, `/mappings`, nil))
}

func (c *APIClient) MappingsCreate(s *Mapping) (string, error) {
	resp, err := easyFetch[RespCreated](c, `POST`, `/mappings`, s)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// func (c *APIClient) MappingsUpdate(id string, s *Mapping) error {
// 	return easyNilFetch(c, `PUT`, `/mappings/` + id, s)
// }

// func (c *APIClient) MappingsDelete(id string) error {
// 	return easyNilFetch(c, `DELETE`, `/mappings/` + id, nil)
// }
