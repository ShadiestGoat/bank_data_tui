package categories

import (
	"fmt"
	"strconv"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/rivo/uniseg"
)

type categoryProxy api.Category

func (c categoryProxy) FilterValue() string {
	return c.Icon + " " + c.Name
}
func (c categoryProxy) GetID() string {
	return c.ID
}
func (c *categoryProxy) SetID(id string) {
	c.ID = id
}

func verifyColor(s string) error {
	if len(s) != 6 {
		return fmt.Errorf("Needs a hex color (no #)")
	}
	if _, err := strconv.ParseUint(s, 16, 64); err != nil {
		return fmt.Errorf("Not a valid color")
	}

	return nil
}

func (c *categoryImpl) NewEditor(w, h int, v *categoryProxy) *editor.Model {
	return editor.New(
		w-listeditor.WIDTH_OFFSET_EDITOR,
		v.ID,
		[]*editor.DataField{
			{
				Title: "Name",
				ID:    "name",
				Value: &v.Name,
				Row: 0,
			},
			{
				Title: "Color",
				ID:    "color",
				Value: &v.Color,
				Row: 1,
				StyleCB: func(v, err string, selected bool, cur lipgloss.Style) lipgloss.Style {
					if !selected || err != "" {
						return cur
					}

					return cur.BorderForeground(lipgloss.Color("#" + v)).BorderStyle(lipgloss.ASCIIBorder())
				},
			},
			{
				Title: "Icon",
				ID:    "icon",
				Value: &v.Icon,
				Row: 2,
			},
		},
		func() (string, error) {
			id, err := c.api.CategoriesCreate(&v.SavableCategory)
			if err != nil {
				return "", err
			}
			return id, nil
		},
		func(id string) error { return c.api.CategoriesUpdate(id, &v.SavableCategory) },
		func(id string) error { return c.api.CategoriesDelete(id) },
		editor.RequireFields(0, 1, 2),
		editor.AddFieldValidator(1, func(s string) error {
			return verifyColor(s)
		}),
		editor.AddFieldValidator(2, func(s string) error {
			if uniseg.GraphemeClusterCount(ansi.Strip(s)) != 1 {
				return fmt.Errorf("Need icon that is 1 in width")
			}

			return nil
		}),
	)
}
