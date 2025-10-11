package editor

import (
	"slices"

	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width int

	id *string

	focusedField   int
	editableFields []*string
	titles         []string
	textFields     []textinput.Model

	popupVisible bool
	popupOnNo    bool

	create func() (string, error)
	update func() error
}

func New(w int, id *string, editableFields []*string, fieldTitles []string, createFunc func() (string, error), updateFunc func() error, mods ...FieldsMod) *Model {
	fields := make([]textinput.Model, len(fieldTitles))
	for i, t := range fieldTitles {
		f := textinput.New()
		f.Prompt = ""
		f.Blur()
		f.Width = 15
		f.Placeholder = t

		fields[i] = f
	}

	for _, m := range mods {
		m(fields)
	}

	for i, f := range editableFields {
		fields[i].SetValue(*f)
	}

	return &Model{
		width:          w,
		id:             id,
		editableFields: editableFields,
		titles:         fieldTitles,
		textFields:     fields,
		create:         createFunc,
		update:         updateFunc,
	}
}

func (c *Model) Init() tea.Cmd {
	cmd := c.textFields[0].Focus()

	return cmd
}

func (c *Model) save() error {
	if c.id == nil {
		id, err := c.create()
		if err != nil {
			return err
		}

		c.id = &id
	} else {
		return c.update()
	}

	return nil
}

func (c *Model) focusField(f int) tea.Cmd {
	if !c.inButtons(c.focusedField) {
		c.textFields[c.focusedField].Blur()
	}
	c.focusedField = f
	if c.inButtons(c.focusedField) {
		return nil
	}

	return c.textFields[c.focusedField].Focus()
}

func (c *Model) incFocusField(d int) tea.Cmd {
	nf := c.focusedField + d

	if nf < 0 {
		// focus on the left button (save) in case of overflow to the minus
		return c.focusField(len(c.textFields))
	} else if nf > len(c.textFields)+1 {
		return c.focusField(0)
	}

	return c.focusField(c.focusedField + d)
}

func (c Model) inButtons(i int) bool {
	return i >= len(c.textFields)
}

func (c *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	batcher := make([]tea.Cmd, 0, len(c.textFields)+1)

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
					batcher = append(batcher, c.focusField(len(c.textFields)-1))
				}
			} else if !c.inButtons(c.focusedField) && (msg.String() == "left" || msg.String() == "right") {
				passToChildren = true
			} else {
				batcher = append(batcher, c.incFocusField(dir))
			}
		case "enter":
			passToChildren = false
			switch c.focusedField {
			case len(c.textFields):
				// save
				if utils.All(slices.Values(c.textFields), func(v textinput.Model) bool { return v.Err == nil }) {
					// TODO: Implement locking mechanism so that this op doesn't block the entire app
					if err := c.save(); err != nil {
						panic(err)
					}
				}
			case len(c.textFields) + 1:
				// reset
				c.focusField(0)
				for i, ogContent := range c.editableFields {
					c.textFields[i].SetValue(*ogContent)
				}
			default:
				batcher = append(batcher, c.focusField(c.focusedField+1))
			}
		}
	}

	if passToChildren {
		for i, f := range c.textFields {
			c.textFields[i], cmd = f.Update(msg)
			batcher = append(batcher, cmd)
		}
	}

	return c, tea.Batch(batcher...)
}

func (c *Model) SetWidth(w int) {
	c.width = w
}

func (c *Model) SetIDPtr(id *string) {
	c.id = id
}
