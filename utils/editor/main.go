package editor

import (
	"fmt"
	"slices"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DataField struct {
	Title  string
	Value  *string
	ID     string

	Row      int
	Col      int
	StyleCB  func(v string, err string, selected bool, cur lipgloss.Style) lipgloss.Style
}

type Model struct {
	width int

	ItemID string

	focusedField int
	dataFields   []*DataField
	inpFields    []textinput.Model
	layout       [][]int

	popupVisible bool
	popupOnNo    bool

	create      func() (string, error)
	update, del func(id string) error
}

func New(w int, id string, dataFields []*DataField, createFunc func() (string, error), updateFunc, delFunc func(id string) error, mods ...FieldsMod) *Model {
	inpFields := make([]textinput.Model, len(dataFields))

	highestRow := 0
	for _, d := range dataFields {
		if d.Row < 0 || d.Col < 0 {
			panic("Row or col can't be <= 0!")
		}

		if d.Row > highestRow {
			highestRow = d.Row
		}
	}

	highestCol := make([]int, highestRow + 1)
	for _, d := range dataFields {
		if d.Col > highestCol[d.Row] {
			highestCol[d.Row] = d.Col
		}
	}

	layout := make([][]int, highestRow + 1)
	for i, v := range highestCol {
		layout[i] = make([]int, v + 1)
	}

	for i, d := range dataFields {
		f := textinput.New()
		f.Prompt = ""
		f.Blur()
		f.Width = 15
		f.Placeholder = d.Title

		inpFields[i] = f

		if layout[d.Row][d.Col] != 0 {
			panic(fmt.Sprintf("Overlap at y=%v, x=%v", d.Row, d.Col))
		}

		// + 1 so that unset detection is simpler :3
		layout[d.Row][d.Col] = i + 1
	}

	for y, r := range layout {
		if len(r) == 0 {
			panic(fmt.Sprintf("Empty row at y=%v", y))
		}
		for x, v := range r {
			if v == 0 {
				panic(fmt.Sprintf("Value not set at y=%v, x=%v", y, x))
			}
			r[x]--
		}
	}

	layout = append(layout, []int{-1, -2, -3})

	for _, m := range mods {
		m(inpFields)
	}

	for i, f := range dataFields {
		inpFields[i].SetValue(*f.Value)
	}

	m := &Model{
		width:      w,
		ItemID:     id,
		dataFields: dataFields,
		inpFields:  inpFields,
		create:     createFunc,
		update:     updateFunc,
		layout:     layout,
		del:        delFunc,
	}

	// For future field growth support
	m.SetWidth(w)

	return m
}

func (c *Model) Init() tea.Cmd {
	cmd := c.inpFields[0].Focus()

	return cmd
}

type ItemNew string
type ItemUpdate string
type ItemDel string

func (c *Model) save() (tea.Msg, error) {
	if c.ItemID == "" {
		id, err := c.create()
		if err != nil {
			return nil, err
		}

		c.ItemID = id
		return ItemNew(id), nil
	}

	err := c.update(c.ItemID)
	if err != nil {
		return nil, err
	}
	return ItemUpdate(c.ItemID), nil
}


func (c *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	batcher := make([]tea.Cmd, 0, len(c.inpFields)+1)

	passToChildren := true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.popupVisible {
			passToChildren = false
		}

		switch msg.String() {
		case "tab", "right", "down", "left", "shift+tab", "up":
			passToChildren = false

			if c.popupVisible {
				c.popupOnNo = !c.popupOnNo
				break
			}

			handled, nf := c.handleNavKey(msg.String())
			if !handled {
				passToChildren = true
			} else {
				batcher = append(batcher, c.focusField(nf))
			}
		case "enter":
			passToChildren = false
			switch c.focusedField {
			case -1:
				// save
				batcher = append(batcher, c.handleSaveEnter())
			case -2:
				// delete
				err := c.del(c.ItemID)
				if err != nil {
					// TODO: Better error handling lmao
					panic("Can't delete: " + err.Error())
				}
				batcher = append(batcher, func() tea.Msg { return ItemDel(c.ItemID) })
			case -3:
				// reset
				c.focusField(0)
				for i, d := range c.dataFields {
					c.inpFields[i].SetValue(*d.Value)
				}
			default:
				batcher = append(batcher, c.focusField(c.focusedField+1))
			}
		}
	case validationErrMsg:
		for _, v := range msg {
			i := slices.IndexFunc(c.dataFields, func(f *DataField) bool { return f.ID == v[0] })
			if i == -1 {
				continue
			}

			if len(c.layout[c.dataFields[i].Row]) != 1 {
				c.inpFields[i].SetValue("")
			}
			c.inpFields[i].Err = APIErr(v[1])
		}
	}

	if passToChildren {
		for i, f := range c.inpFields {
			c.inpFields[i], cmd = f.Update(msg)
			batcher = append(batcher, cmd)
		}
	}

	return c, tea.Batch(batcher...)
}

func (c *Model) SetWidth(w int) {
	c.width = w
}

type validationErrMsg [][2]string

func (c *Model) handleSaveEnter() tea.Cmd {
	if utils.Any(slices.Values(c.inpFields), func(v textinput.Model) bool { return v.Err != nil }) {
		return nil
	}

	for i, f := range c.inpFields {
		*c.dataFields[i].Value = f.Value()
	}

	return func() tea.Msg {
		msg, err := c.save()
		if err == nil {
			return msg
		}

		if e, ok := err.(*api.ValidationErr); !ok {
			panic(err)
		} else {
			return validationErrMsg(e.Details)
		}
	}
}
