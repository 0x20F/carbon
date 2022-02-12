package logger

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Alignment int

const (
	Left   Alignment = iota
	Center Alignment = iota
	Right  Alignment = iota
)

// Table structure to allow for easy printing and formatting
// of data in a table.
//
// Built in accounting for any weird color characters that might
// show up in the data that needs to be printed, as well as
// scalable column width based on automatic detection of the
// longest value in a given column.
//
// Default alignment is Left alignment.
type Table struct {
	ColumnCount int        // How many columns the table will have
	Paddings    []int      // The widths of each column
	Columns     [][]string // The data to be printed
}

func NewTable(columnCount int) *Table {
	columns := make([][]string, columnCount)

	table := &Table{
		ColumnCount: columnCount,
		Paddings:    make([]int, columnCount),
		Columns:     columns,
	}

	return table
}

func (t *Table) Header(data ...string) {
	t.addRow(data)

	spacer := make([]string, t.ColumnCount)
	t.addRow(spacer)
}

func (t *Table) AddRow(data ...string) {
	t.addRow(data)
}

func (t *Table) Display() {
	rowCount := len(t.Columns[0])

	for i := 0; i < rowCount; i++ {
		current := make([]string, t.ColumnCount)

		for j := 0; j < t.ColumnCount; j++ {
			value := t.Columns[j][i]
			formatted := formatString(value, t.Paddings[j])
			current = append(current, formatted)
		}

		fmt.Println(strings.Join(current, ""))
	}
}

func (t *Table) addRow(row []string) {
	for i := 0; i < t.ColumnCount; i++ {
		item := row[i]

		// If the current element is longer than the current
		// padding, update the padding
		if len(item) > t.Paddings[i] {
			t.Paddings[i] = len(item)
		}

		// Add it to the correct column
		t.Columns[i] = append(t.Columns[i], item)
	}
}

func formatString(str string, pad int) string {
	// Find all ANSI escape codes in the string
	ansiCodes := regexp.MustCompile("\x1b\\[[0-9;]*[mK]")
	ansiCount := ansiCodes.FindAllStringIndex(str, -1)

	width := pad + len(ansiCount) + 2 // the 2 is to have an extra space on either side

	// Pad the string
	final := fmt.Sprintf("%-"+strconv.Itoa(width)+"s", str)

	return final
}
