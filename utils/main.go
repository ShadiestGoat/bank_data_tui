package utils

import (
	"iter"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func JoinHorizontal2(w int, a, b string) string {
	res, _ := JoinHorizontalWithSpacer(w, 1, a, b)
	return res
}

func JoinHorizontalWithSpacer(w, spacerIndex int, str ...string) (string, []int) {
	offsets := make([]int, len(str))
	widths := 0
	for _, s := range str {
		widths += lipgloss.Width(s)
	}

	str = slices.Insert(str, spacerIndex, strings.Repeat(" ", w-widths))
	off := 0
	for i, v := range str {
		if spacerIndex == i {
			off += lipgloss.Width(v)
			continue
		} else if i > spacerIndex {
			i--
		}

		offsets[i] = off
		off += lipgloss.Width(v)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, str...), offsets
}

// Returns a rendered str, and offsets from the left
func JoinHorizontalEqualSpread(w int, str ...string) (string, []int) {
	if len(str) == 0 {
		return "", nil
	} else if len(str) == 1 {
		usedW := lipgloss.Width(str[0])
		leftover := w - usedW

		return lipgloss.NewStyle().PaddingLeft(leftover/2).PaddingRight(leftover - leftover/2).Render(str[0]), []int{leftover/2}
	}

	widths := 0
	for _, s := range str {
		widths += lipgloss.Width(s)
	}

	leftover := w - widths
	perSlice := leftover / (len(str) - 1)
	if perSlice < 0 {
		return "", nil
	}

	addExtraSpaceEvery := (perSlice % (len(str) - 1)) + 1

	res := make([]string, len(str)+len(str)-1)
	curOff := 0
	offsets := make([]int, len(str))
	for i := range res {
		if i%2 == 0 {
			res[i] = str[i/2]
			offsets[i/2] = curOff
			curOff += lipgloss.Width(str[i/2])
		} else {
			res[i] = strings.Repeat(" ", perSlice)
			curOff += perSlice
			if addExtraSpaceEvery != 1 && i%addExtraSpaceEvery == 0 {
				curOff++
				res[i] += " "
			}
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, res...), offsets
}

// Reports if any in sl is true
func Any[T any](sl iter.Seq[T], cond func(T) bool) bool {
	for v := range sl {
		if cond(v) {
			return true
		}
	}

	return false
}

// Reports if all in sl is true
func All[T any](sl iter.Seq[T], cond func(T) bool) bool {
	for v := range sl {
		if !cond(v) {
			return false
		}
	}

	return true
}

func Overflow(str string, maxWidth int) string {
	if lipgloss.Width(str) <= maxWidth {
		return str
	}

	lines := strings.Split(str, "\n")
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], maxWidth, "…")
	}

	return strings.Join(lines, "\n")
}

type ResizeMessage struct {
	W, H int
}

type Screen interface {
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() (string, *tea.Cursor)
	Init() tea.Cmd
}

type MsgGoToHome struct {}
func GoToHome() tea.Msg { return MsgGoToHome{} }