package mappings

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/bubbles/textinput"
)

type mappingProxy api.Mapping

func (m mappingProxy) FilterValue() string {
	return m.Name
}
func (m mappingProxy) GetID() string {
	return m.ID
}
func (m *mappingProxy) SetID(id string) {
	m.ID = id
}

func (c *mappingImpl) NewEditor(w, h int, v *mappingProxy) *editor.Model {
	return editor.New(
		w-listeditor.WIDTH_OFFSET_EDITOR,
		v.ID,
		[]*editor.DataField{
			{
				Title: "Name",
				ID:    "name",
				Value: &v.Name,
				Row:   0,
				Flex:  true,
			},
			{
				Title: "Priority",
				ID:    "priority",
				GetValue: func() string {
					if v.Priority == 0 {
						return ""
					}
					return strconv.Itoa(v.Priority)
				},
				SetValue: func(raw string) {
					parsed, _ := strconv.Atoi(raw)
					// Validation handles err handling
					v.Priority = parsed
				},
				Row: 0,
				Col: 1,
			},
			{
				Title: "Match Description Regex",
				ID:    "inpText",
				Value: &v.InpText,
				Row:   1,
				Flex:  true,
			},
			{
				Title: "Match Amount",
				ID:    "inpAmt",
				GetValue: func() string {
					if v.InpAmt == nil {
						return ""
					}

					return strconv.FormatFloat(*v.InpAmt, 'e', 2, 64)
				},
				SetValue: func(raw string) {
					if raw == "" {
						v.InpAmt = nil
						return
					}

					parsed, _ := strconv.ParseFloat(raw, 64)
					// Validation handles err handling
					v.InpAmt = &parsed
				},
				Row: 1,
				Col: 1,
			},
			{
				Title: "Resulting Name",
				ID:    "resName",
				Value: &v.ResName,
				Row:   2,
				Flex:  true,
			},
			{
				Title: "Resulting Category",
				ID:    "resCategory",
				GetValue: func() string {
					if v.ResCategoryID == "" {
						return ""
					}
					i := slices.IndexFunc(c.cache.Categories, func(c *api.Category) bool {
						return c.ID == v.ResCategoryID
					})
					if i == -1 {
						return ""
					}

					return c.cache.Categories[i].Name
				},
				SetValue: func(raw string) {
					if raw == "" {
						v.ResCategoryID = ""
						return
					}

					for _, c := range c.cache.Categories {
						if strings.EqualFold(raw, c.Name) {
							v.ResCategoryID = c.ID
							return
						}
					}

					v.ResCategoryID = ""
				},
				Row:  2,
				Col:  1,
				Flex: true,
			},
		},
		func(alt bool) (string, error) {
			id, err := c.api.MappingsCreate((*api.Mapping)(v), alt)
			if err != nil {
				return "", err
			}
			return id, nil
		},
		func(alt bool, id string) error {
			err := c.api.MappingsUpdate(id, (*api.Mapping)(v), alt)
			if err != nil {
				return err
			}
			return nil
		},
		func(alt bool, id string) error {
			err := c.api.MappingsDelete(id, alt)
			if err != nil {
				return err
			}
			return nil
		},
		editor.RequireFields(0),
		editor.AddIntValidator(1),
		editor.AddFloatValidator(3),
		editor.AddFieldValidator(2, func(s string) error {
			if s == "" {
				return nil
			}
			_, err := regexp.CompilePOSIX(s)
			if err != nil {
				return fmt.Errorf("Must be a valid (posix) regex")
			}

			return nil
		}),
		editor.AddFieldValidator(5, func(s string) error {
			if s == "" {
				return nil
			}

			for _, c := range c.cache.Categories {
				if strings.EqualFold(s, c.Name) {
					return nil
				}
			}

			return fmt.Errorf("Must be a valid category")
		}),
		func(fields []*textinput.Model) {
			fields[5].ShowSuggestions = true
			c.categoryField = fields[5]
			c.resetSuggestions()
		},
		editor.AddOneOfRequirement("matcher", 2, 3),
		editor.AddOneOfRequirement("result", 4, 5),
	)
}
