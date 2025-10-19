package editor

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DataField struct {
	Title string
	ID    string

	Row int
	Col int

	Value    *string
	GetValue func() string
	SetValue func(v string)

	Flex    bool
	StyleCB func(v string, err error, selected bool, cur lipgloss.Style) lipgloss.Style
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
		if d.Value == nil && (d.GetValue == nil || d.SetValue == nil) {
			panic("Data Field must have at least 1 field get/set method")
		}

		if d.Row > highestRow {
			highestRow = d.Row
		}
	}

	highestCol := make([]int, highestRow+1)
	for _, d := range dataFields {
		if d.Col > highestCol[d.Row] {
			highestCol[d.Row] = d.Col
		}
	}

	layout := make([][]int, highestRow+1)
	for i, v := range highestCol {
		layout[i] = make([]int, v+1)
	}

	for i, d := range dataFields {
		f := textinput.New()
		f.Prompt = ""
		f.Blur()
		f.Width = 15
		f.Placeholder = d.Title
		f.KeyMap.NextSuggestion.SetKeys("ctrl+n")
		f.KeyMap.PrevSuggestion.SetKeys("ctrl+p")

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

	for _, m := range mods {
		m(inpFields)
	}

	for i, f := range dataFields {
		if f.Value == nil {
			inpFields[i].SetValue(f.GetValue())
		} else {
			inpFields[i].SetValue(*f.Value)
		}
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

	m.SetWidth(w)
	m.resetButtonLayout()

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
		c.resetButtonLayout()
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
			case BTN_SAVE:
				// save
				batcher = append(batcher, c.handleSaveEnter())
			case BTN_DEL:
				// delete
				err := c.del(c.ItemID)
				if err != nil {
					// TODO: Better error handling lmao
					panic("Can't delete: " + err.Error())
				}
				batcher = append(batcher, func() tea.Msg { return ItemDel(c.ItemID) })
			case BTN_RESET:
				// reset
				c.focusField(c.layout[0][0])
				for i, d := range c.dataFields {
					if d.Value == nil {
						c.inpFields[i].SetValue(d.GetValue())
					} else {
						c.inpFields[i].SetValue(*d.Value)
					}
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
		if _, ok := msg.(tea.KeyMsg); ok {
			for i, f := range c.inpFields {
				// re-validate cause some validators need to be triggered external events
				if errors.Is(f.Err, APIErr("")) {
					continue
				}
				c.inpFields[i].Err = f.Validate(f.Value())
			}
		}
	}

	return c, tea.Batch(batcher...)
}

func (c *Model) SetWidth(w int) {
	c.width = w
	log.Println("WIDTH", w)

	for _, row := range c.layout {
		if row[0] < 0 {
			continue
		}

		// -len(row) + 1 = space between each item
		// -len(row) = each item has an extra space aside from width
		// fuck you textinput component >:(
		availWidth := w - len(row) + 1 - len(row)
		flexers := 0

		for _, i := range row {
			f := c.dataFields[i]
			availWidth -= extraFieldLength(&c.inpFields[i], f)

			if f.Flex {
				flexers++
			} else {
				availWidth -= c.inpFields[i].Width
			}
		}

		if flexers == 0 {
			continue
		}

		extraSpaceEvery := availWidth % flexers
		spaceBuf := 0

		for _, i := range row {
			if !c.dataFields[i].Flex {
				continue
			}
			c.inpFields[i].Width = availWidth / flexers
			if extraSpaceEvery != 0 && spaceBuf == extraSpaceEvery {
				c.inpFields[i].Width += 1
				spaceBuf = 0
			} else {
				spaceBuf++
			}
		}
	}
}

type validationErrMsg [][2]string

func (c *Model) handleSaveEnter() tea.Cmd {
	if utils.Any(slices.Values(c.inpFields), func(v textinput.Model) bool { return v.Err != nil }) {
		return nil
	}

	for i, f := range c.inpFields {
		d := c.dataFields[i]
		if d.Value == nil {
			d.SetValue(f.Value())
		} else {
			*d.Value = f.Value()
		}
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

const (
	BTN_SAVE  = -1
	BTN_DEL   = -2
	BTN_RESET = -3
)

func (c *Model) resetButtonLayout() {
	y := len(c.layout) - 1
	var l []int
	if c.ItemID == "" {
		l = []int{BTN_SAVE, BTN_RESET}
	} else {
		l = []int{BTN_SAVE, BTN_DEL, BTN_RESET}
	}

	if c.layout[y][0] < 0 {
		c.layout[y] = l
	} else {
		c.layout = append(c.layout, l)
	}
}
