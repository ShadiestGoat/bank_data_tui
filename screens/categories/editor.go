package categories

import (
	"fmt"
	"strconv"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils/editor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/rivo/uniseg"
)

func verifyColor(s string) error {
	if len(s) != 6 {
		return fmt.Errorf("Needs a hex color (no #)")
	}
	if _, err := strconv.ParseUint(s, 16, 64); err != nil {
		return fmt.Errorf("Not a valid color")
	}

	return nil
}

func newCategoryEditor(c *api.APIClient, w int, cat *api.Category) *editor.Model {
	return editor.New(
		w-WIDTH_OFFSET_EDITOR,
		cat.ID,
		[]*editor.DataField{
			{
				Title: "Name",
				ID:    "name",
				Value: &cat.Name,
			},
			{
				Title: "Color",
				ID:    "color",
				Value: &cat.Color,
			},
			{
				Title: "Icon",
				ID:    "icon",
				Value: &cat.Icon,
			},
		},
		func() (string, error) {
			id, err := c.CategoriesCreate(&cat.SavableCategory)
			if err != nil {
				return "", err
			}
			return id, nil
		},
		func(id string) error { return c.CategoriesUpdate(id, &cat.SavableCategory) },
		func(id string) error { return c.CategoriesDelete(id) },
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

func (m *Model) resetCategoryEditor() tea.Cmd {
	m.editor = newCategoryEditor(m.c, m.w, m.curItem)
	return m.editor.Init()
}
