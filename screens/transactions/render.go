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

//  1     60%     40%    10     8
// ICON | NAME | DESC | DATE | AMT
//

const COL_SPLIT = "│"

func (m *Model) cols() []int {
	colCunt := 4

	icon := 2
	amt := 8
	date := 10

	// space padding on either side + content
	colTotalWidth := lipgloss.Width(COL_SPLIT)*colCunt + colCunt * 2
	leftover := m.w - icon - amt - date - colTotalWidth

	nameLen := int(float64(leftover) * 0.6)

	return []int{icon, nameLen, leftover - nameLen, date, amt}
}

func (m Model) renderRow(t *api.Transaction, selected bool) string {
	rowStyle := lipgloss.NewStyle()
	if selected {
		rowStyle = rowStyle.Background(styles.COLOR_MAIN)
	}

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

	str[2] = t.Desc
	str[3] = t.AuthedAt.Format("02/01/2006")

	if t.ResolvedName != nil {
		str[1] = *t.ResolvedName
		str[2] = lipgloss.NewStyle().Faint(true).Render(str[2])
	}

	str[4] = strconv.FormatFloat(t.Amount, 'f', 2, 64)

	for i, w := range cols {
		base := rowStyle
		if i == 0 {
			base = lipgloss.NewStyle()
		}

		str[i] = base.Width(w).Render(
			utils.Overflow(str[i], w),
		)
	}

	if cat != nil {
		str[0] = lipgloss.NewStyle().Background(
			lipgloss.Color("#" + cat.Color),
		).Foreground(lipgloss.AdaptiveColor{"#ffffff", "#000000"}).Render(str[0] + " ")
	}

	colSplitter := rowStyle.Render(" " + COL_SPLIT + " ")

	// str[0] alr has a space in it
	return str[0] + COL_SPLIT + rowStyle.Render(
		" "+strings.Join(str[1:], colSplitter),
	)
}

func (m Model) View() string {
	if m.h == 0 || len(m.items) == 0 {
		return ""
	}

	items := m.items[m.viewportOff:]
	// 1 line empty at the bottom
	items = items[:min(len(items), m.h-1)]

	if len(items) == 0 {
		return "No Items here!"
	}

	rows := ""
	for i, v := range items {
		rows += m.renderRow(v, m.selected == m.viewportOff+i) + "\n"
	}

	lastRowItems := []string{"Total Transactions: " + strconv.Itoa(len(m.items))}
	if m.nextPageLoading {
		loading := m.loader.View()
		lastRowItems = append(
			lastRowItems,
			loading+lipgloss.NewStyle().Faint(true).Render(" Loading"),
		)
	}

	rows = rows[:len(rows)-1]

	if m.hasHitLastPage && m.h-len(items) > 3 {
		rows += "\n\n\n" + lipgloss.PlaceHorizontal(m.w, lipgloss.Center, "No More Transactions!")
	}

	return rows + strings.Repeat("\n", m.h-strings.Count(rows, "\n")-1) + utils.JoinHorizontalWithSpacer(m.w, 1, lastRowItems...)
}
