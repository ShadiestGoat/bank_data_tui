package utils

import (
	"iter"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func JoinHorizontal2(w int, a, b string) string {
	return JoinHorizontalWithSpacer(w, 1, a, b)
}

func JoinHorizontalWithSpacer(w, spacerIndex int, str ...string) string {
	widths := 0
	for _, s := range str {
		widths += lipgloss.Width(s)
	}

	str = slices.Insert(str, spacerIndex, strings.Repeat(" ", w-widths))

	return lipgloss.JoinHorizontal(lipgloss.Center, str...)
}

func JoinHorizontalEqualSpread(w int, str ...string) string {
	if len(str) == 0 {
		return ""
	} else if len(str) == 1 {
		return lipgloss.PlaceHorizontal(w, lipgloss.Center, str[0])
	}

	widths := 0
	for _, s := range str {
		widths += lipgloss.Width(s)
	}

	leftover := w - widths
	perSlice := leftover/(len(str) - 1)
	if perSlice < 0 {
		return ""
	}

	addExtraSpaceEvery := (perSlice % (len(str) - 1)) + 1

	res := make([]string, len(str) + len(str) - 1)
	for i := range res {
		if i % 2 == 0 {
			res[i] = str[i/2]
		} else {
			res[i] = strings.Repeat(" ", perSlice)
			if addExtraSpaceEvery != 1 && i % addExtraSpaceEvery == 0 {
				res[i] += " "
			}
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, res...)
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
