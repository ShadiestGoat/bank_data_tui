package categories

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils/editor"
	"github.com/bank_data_tui/utils/listeditor"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
			},
			{
				Title: "Color",
				ID:    "color",
				Value: &v.Color,
			},
			{
				Title: "Icon",
				ID:    "icon",
				Value: &v.Icon,
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

type categoryDelegate struct{}

func (c categoryDelegate) Spacing() int { return 1 }
func (c categoryDelegate) Height() int  { return 1 }

func (c categoryDelegate) Render(w io.Writer, m list.Model, i int, v list.Item) {
	style := lipgloss.NewStyle().Foreground(styles.COLOR_MAIN)
	if m.GlobalIndex() == i {
		style = style.Underline(true)
	}

	txt, ok := v.(listeditor.NewItem)
	if ok {
		w.Write(
			[]byte(" " + style.Render(string(txt))),
		)

		return
	}

	cat := v.(*categoryProxy)
	if err := verifyColor(cat.Color); err == nil {
		style = style.Foreground(lipgloss.Color("#" + cat.Color))
	}

	w.Write(
		[]byte(" " + style.Render("["+cat.Icon+"] "+cat.Name)),
	)
}

func (c categoryDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
