package mappings

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

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
				Value: &v.InpText,
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
					i := slices.IndexFunc(c.fetchedCategories, func(c *api.Category) bool {
						return c.ID == v.ResCategoryID
					})
					if i == -1 {
						return ""
					}

					return c.fetchedCategories[i].Name
				},
				SetValue: func(raw string) {
					if raw == "" {
						v.ResCategoryID = ""
						return
					}

					for _, c := range c.fetchedCategories {
						if c.Name == raw {
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
		func() (string, error) {
			id, err := c.api.MappingsCreate((*api.Mapping)(v))
			if err != nil {
				return "", err
			}
			return id, nil
		},
		func(id string) error { panic("unimplemented") },
		func(id string) error { panic("unimplemented") },
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

			for _, c := range c.fetchedCategories {
				if s == c.Name {
					return nil
				}
			}

			return fmt.Errorf("Must be a valid category")
		}),
		func(fields []textinput.Model) {
			fields[5].ShowSuggestions = true
			c.categoryField = &fields[5]
		},
		editor.AddOneOfRequirement("matcher", 2, 3),
		editor.AddOneOfRequirement("result", 4, 5),
	)
}
