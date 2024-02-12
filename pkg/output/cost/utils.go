package cost

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/fatih/color"
	"github.com/leekchan/accounting"
	"sort"
	"strconv"
)

var bold = color.New(color.Bold)

func sortRows(rows []table.Row) []table.Row {
	sort.Slice(rows, func(i, j int) bool {
		numI, _ := strconv.ParseFloat(rows[i][len(rows[i])-1], 64)
		numJ, _ := strconv.ParseFloat(rows[j][len(rows[j])-1], 64)
		return numI > numJ
	})

	return rows
}

func makeNumbersAccounting(rows []table.Row) []table.Row {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	for _, row := range rows {
		costFloat, _ := strconv.ParseFloat(row[len(row)-1], 64)
		row[len(row)-1] = ac.FormatMoney(costFloat)
	}
	return rows
}
