package utils

import (
	"iter"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func JoinHorizontal2(w int, a, b string) string {
	return JoinHorizontalSpread(w, 1, a, b)
}

func JoinHorizontalSpread(w, spacerIndex int, str ...string) string {
	widths := 0
	for _, s := range str {
		widths += lipgloss.Width(s)
	}

	str = slices.Insert(str, spacerIndex, strings.Repeat(" ", w - widths))

	return lipgloss.JoinHorizontal(lipgloss.Center, str...)
}

// Reports if any in sl is true
func Any[T any](sl iter.Seq[T], cond func (T) bool) bool {
	for v := range sl {
		if cond(v) {
			return true
		}
	}

	return false
}

// Reports if all in sl is true
func All[T any](sl iter.Seq[T], cond func (T) bool) bool {
	for v := range sl {
		if !cond(v) {
			return false
		}
	}

	return true
}
