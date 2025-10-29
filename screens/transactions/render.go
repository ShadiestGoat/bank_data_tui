package transactions

import (
	"slices"
	"strconv"
	"strings"

	"github.com/bank_data_tui/api"
	"github.com/bank_data_tui/styles"
	"github.com/bank_data_tui/utils"
	"github.com/charmbracelet/lipgloss"
)

//  1     60%     40%     8
// ICON | NAME | DESC | AMT
//

const COL_SPLIT = " | "

func (m *Model) cols() []int {
	icon := 2
	amt := 8

	leftover := m.w - icon - amt - len(COL_SPLIT)*3

	nameLen := int(float64(leftover) * 0.6)

	return []int{icon, nameLen, leftover - nameLen, amt}
}

func (m *Model) renderRow(t *api.Transaction, selected bool) string {
	cols := m.cols()
	str := make([]string, len(cols))
	var cat *api.Category
	if t.ResolvedCategoryID != nil {
		i := slices.IndexFunc(m.cache.Categories, func(c *api.Category) bool {
			return c.ID == *t.ResolvedCategoryID
		})
		if i != -1 {
			cat = m.cache.Categories[i]
			str[0] = m.cache.Categories[i].Icon
		}
	}

	if len(cols) == 3 {
		if t.ResolvedName != nil {
			str[1] = *t.ResolvedName
		} else {
			str[1] = lipgloss.NewStyle().Faint(true).Italic(true).Render(t.Desc)
		}
	} else {
		str[2] = t.Desc
		if t.ResolvedName != nil {
			str[1] = *t.ResolvedName
			str[2] = lipgloss.NewStyle().Faint(true).Italic(true).Render(str[2])
		}
	}

	str[len(str)-1] = strconv.FormatFloat(t.Amount, 'f', 2, 64)

	for i, w := range cols {
		str[i] = lipgloss.NewStyle().Width(w).Render(
			utils.Overflow(str[i], w),
		)
	}

	if cat != nil {
		str[0] = lipgloss.NewStyle().Background(
			lipgloss.Color("#" + cat.Color),
		).Foreground(lipgloss.AdaptiveColor{"#ffffff", "#000000"}).Render(str[0] + " ")
	}

	rowStyle := lipgloss.NewStyle()
	if selected {
		rowStyle = rowStyle.Background(styles.COLOR_MAIN)
	}

	return str[0] + COL_SPLIT[1:len(COL_SPLIT)-1] + rowStyle.Render(
		" "+strings.Join(str[1:], COL_SPLIT),
	)
}

func (m *Model) View() string {
	if m.h == 0 || len(m.items) == 0 {
		return ""
	}

	items := m.items[m.viewportOff:]
	// 1 line empty at the bottom
	items = items[:min(len(items), m.h - 1)]

	if len(items) == 0 {
		return "No Items here!"
	}

	rows := ""
	for i, v := range items {
		rows += m.renderRow(v, m.selected == m.viewportOff + i) + "\n"
	}


	return rows[:len(rows)-1]
}
