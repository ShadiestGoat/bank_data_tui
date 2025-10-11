package main

import (
	"strings"

	"github.com/bank_data_tui/styles"
	"github.com/charmbracelet/lipgloss"
)

var (
	STYLE_HEADER_TEXT     = lipgloss.NewStyle().Foreground(styles.COLOR_MAIN).Margin(1)
	STYLE_HEADER_SELECTED = STYLE_HEADER_TEXT.Bold(true).Underline(true)
	STYLE_HEADER          = lipgloss.NewStyle().Border(lipgloss.DoubleBorder(), false, false, true, false).Margin(0, 0, 2, 0).BorderForeground(styles.COLOR_MAIN)
)

const (
	// 2 (margin bottom) + 1 * 2 (padding top & bot) + 1 line of border + 1 line of text
	HEADER_HEIGHT = 6
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

func (m mainApp) View() string {
	box := lipgloss.NewStyle().Width(m.width).Height(m.height).AlignHorizontal(lipgloss.Center)

	if m.width == 0 || m.height == 0 {
		return ""
	} else if m.width < 30 || m.height < 20 {
		return box.AlignVertical(lipgloss.Center).Render("Too small")
	}

	s := m.screenImp.View()
	if m.curFocusedScreen == S_LOGIN {
		return box.AlignVertical(lipgloss.Center).Render(s)
	}

	h := m.renderHeader()

	return box.AlignHorizontal(lipgloss.Left).Render(lipgloss.JoinVertical(lipgloss.Center, h, s))
}
