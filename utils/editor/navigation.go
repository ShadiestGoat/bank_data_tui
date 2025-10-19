package editor

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (c *Model) focusField(f int) tea.Cmd {
	oldPos := 0
	if !c.inButtons(c.focusedField) {
		c.inpFields[c.focusedField].Blur()
		oldPos = c.inpFields[c.focusedField].Position()
	}

	c.focusedField = f
	if c.inButtons(c.focusedField) {
		return nil
	}

	cmd := c.inpFields[c.focusedField].Focus()
	c.inpFields[c.focusedField].SetCursor(oldPos)

	return cmd
}

func (c Model) inButtons(i int) bool {
	return i < 0
}

// returns "handled", "new focus id"
func (c Model) handleNavKey(key string) (bool, int) {
	switch key {
	case "tab":
		return true, c.navKeyHorizontal(1)
	case "shift+tab":
		return true, c.navKeyHorizontal(-1)
	case "down":
		return true, c.navKeyVertical(1)
	case "up":
		return true, c.navKeyVertical(-1)
	case "right":
		return c.navKeyHorizontalTextConflict(1)
	case "left":
		return c.navKeyHorizontalTextConflict(-1)
	}

	return false, 0
}

func (c Model) rowColForIndex(i int) (int, int) {
	if c.inButtons(i) {
		row := len(c.layout) - 1
		return row, slices.Index(c.layout[row], i)
	}

	v := c.dataFields[i]
	return v.Row, v.Col
}

func (c Model) navKeyHorizontalTextConflict(dir int) (bool, int) {
	if c.inButtons(c.focusedField) {
		return true, c.navKeyHorizontal(dir)
	}
	np := c.inpFields[c.focusedField].Position() + dir
	if np < 0 || np > lipgloss.Width(c.inpFields[c.focusedField].Value()) {
		return true, c.navKeyHorizontal(dir)
	}

	return false, c.focusedField
}

func (c Model) navKeyHorizontal(dir int) int {
	cy, cx := c.rowColForIndex(c.focusedField)

	// handle horizontal first
	if cx + dir >= 0 && cx + dir < len(c.layout[cy]) {
		return c.layout[cy][cx + dir]
	}

	return c.navKeyVertical(dir)
}

func (c Model) navKeyVertical(dir int) int {
	cy, cx := c.rowColForIndex(c.focusedField)

	// Out of bounds first
	if cy + dir < 0 {
		return c.layout[len(c.layout) - 1][0]
	} else if cy + dir >= len(c.layout) {
		return c.layout[0][0]
	}

	curPerc := float64(cx)/float64(len(c.layout[cy]))
	nextLen := float64(len(c.layout[cy + dir]))
	for i := range c.layout[cy + dir] {
		if curPerc >= float64(i)/nextLen {
			return c.layout[cy + dir][max(i - 1, 0)]
		}
	}

	panic("Meow")

	return 0
}
