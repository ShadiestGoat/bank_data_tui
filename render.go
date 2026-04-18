package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/bank_data_tui/styles"
)

var (
	STYLE_HEADER_TEXT     = lipgloss.NewStyle().Foreground(styles.COLOR_MAIN).Margin(1)
	STYLE_HEADER_SELECTED = STYLE_HEADER_TEXT.Bold(true).Underline(true)
	STYLE_HEADER          = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), false, false, true, false).Margin(0, 0, 1, 0).BorderForeground(styles.COLOR_MAIN)
)

const (
	// 1 (margin bottom) + 1 * 2 (padding top & bot) + 1 line of border + 1 line of text
	HEADER_HEIGHT = 1 + 1*2 + 1 + 1
)

var HEADER_SCREENS = []struct {
	s Screen
	t string
}{
	{S_TRANS, "Transactions"},
	{S_MAPPINGS, "Mappings"},
	{S_CATEGORIES, "Categories"},
	{S_UPLOAD, "Upload"},
}

func (m mainApp) renderHeader() string {
	r := []string{}
	for _, h := range HEADER_SCREENS {
		if h.s == m.curFocusedScreen {
			r = append(r, STYLE_HEADER_SELECTED.Render(h.t))
		} else {
			r = append(r, STYLE_HEADER_TEXT.Render(h.t))
		}
	}

	left := r[0]
	right := lipgloss.JoinHorizontal(lipgloss.Top, r[1:]...)

	spacer := strings.Repeat(" ", m.width-lipgloss.Width(right)-lipgloss.Width(left))

	return STYLE_HEADER.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, spacer, right))
}

func (m mainApp) renderTooSmall() string {
	box := lipgloss.NewStyle().Width(m.width).Height(m.height).Align(lipgloss.Center, lipgloss.Center)
	return box.Render("Too Small")
}

func (m mainApp) View() (v tea.View) {
	v.AltScreen = true
	v.WindowTitle = "Bank Data"
	v.Cursor = nil

	if m.width == 0 || m.height == 0 {
		return
	} else if m.width < 50 || m.height < 20 {
		v.SetContent(m.renderTooSmall())
		return
	}

	s, c := m.screenImp.View()
	v.Cursor = c

	w, h := lipgloss.Width(s), lipgloss.Height(s)
	if h > (m.height - HEADER_HEIGHT) || w > m.width {
		v.Cursor = nil
		v.SetContent(m.renderTooSmall())
		return 
	}

	if m.curFocusedScreen == S_LOGIN {
		padTop := (m.height - h)/2
		padLeft := (m.width - w)/2

		if c != nil {
			c.Y += padTop
			c.X += padLeft
		}

		v.SetContent(lipgloss.NewStyle().Padding(padTop, 0, 0, padLeft).Render(s))
		return
	}

	header := m.renderHeader()
	padTop := (m.height - h - HEADER_HEIGHT)/2
	padLeft := (m.width - w)/2

	if c != nil {
		c.Y += padTop + HEADER_HEIGHT
		c.X += padLeft
	}

	v.SetContent(header + "\n" + lipgloss.NewStyle().Padding(padTop, 0, 0, padLeft).Render(s))

	return v
}
