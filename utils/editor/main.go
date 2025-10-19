package editor

import (
	"slices"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DataField struct {
	Title  string
	Value  *string
	ID     string
	apiErr string
}

type Model struct {
	width int

	ItemID string

	focusedField int
	dataFields   []*DataField
	inpFields    []textinput.Model

	popupVisible bool
	popupOnNo    bool

	create      func() (string, error)
	update, del func(id string) error
}

func New(w int, id string, dataFields []*DataField, createFunc func() (string, error), updateFunc, delFunc func(id string) error, mods ...FieldsMod) *Model {
	inpFields := make([]textinput.Model, len(dataFields))
	for i, d := range dataFields {
		f := textinput.New()
		f.Prompt = ""
		f.Blur()
		f.Width = 15
		f.Placeholder = d.Title

		inpFields[i] = f
	}

	for _, m := range mods {
		m(inpFields)
	}

	for i, f := range dataFields {
		inpFields[i].SetValue(*f.Value)
	}

	return &Model{
		width:      w,
		ItemID:     id,
		dataFields: dataFields,
		inpFields:  inpFields,
		create:     createFunc,
		update:     updateFunc,
		del: delFunc,
	}
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

func (c *Model) focusField(f int) tea.Cmd {
	if !c.inButtons(c.focusedField) {
		c.inpFields[c.focusedField].Blur()
	}
	c.focusedField = f
	if c.inButtons(c.focusedField) {
		return nil
	}

	return c.inpFields[c.focusedField].Focus()
}

func (c *Model) incFocusField(d int) tea.Cmd {
	nf := c.focusedField + d

	if nf < 0 {
		// focus on the left button (save) in case of overflow to the minus
		return c.focusField(len(c.inpFields))
	} else if nf > len(c.inpFields)+2 {
		return c.focusField(0)
	}

	return c.focusField(c.focusedField + d)
}

func (c Model) inButtons(i int) bool {
	return i >= len(c.inpFields)
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

			var dir int
			switch msg.String() {
			case "right", "tab", "down":
				dir = 1
			case "left", "shift+tab", "up":
				dir = -1
			}

			if c.inButtons(c.focusedField) && (msg.String() == "down" || msg.String() == "up") {
				// in buttons, a up/down key should yield "get out of the button row"
				if msg.String() == "down" {
					batcher = append(batcher, c.focusField(0))
				} else {
					batcher = append(batcher, c.focusField(len(c.inpFields)-1))
				}
			} else if !c.inButtons(c.focusedField) && (msg.String() == "left" || msg.String() == "right") {
				passToChildren = true
			} else {
				batcher = append(batcher, c.incFocusField(dir))
			}
		case "enter":
			passToChildren = false
			switch c.focusedField {
			case len(c.inpFields):
				// save
				batcher = append(batcher, c.handleSaveEnter())
			case len(c.inpFields) + 1:
				// delete
				err := c.del(c.ItemID)
				if err != nil {
					// TODO: Better error handling lmao
					panic("Can't delete: " + err.Error())
				}
				batcher = append(batcher, func() tea.Msg { return ItemDel(c.ItemID) })
			case len(c.inpFields) + 2:
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

			c.dataFields[i].apiErr = v[1]
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

// TODO: Implement locking mechanism so that this op doesn't block the entire app
func (c *Model) handleSaveEnter() tea.Cmd {
	if utils.Any(slices.Values(c.inpFields), func(v textinput.Model) bool { return v.Err != nil }) {
		return nil
	}

	for i, f := range c.inpFields {
		*c.dataFields[i].Value = f.Value()
		c.dataFields[i].apiErr = ""
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
