package printer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pborman/ansi"
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

// Creates a new table with the given column count.
// This needs to be defined since other internal measurements
// depend on it.
func NewTable(columnCount int) *Table {
	columns := make([][]string, columnCount)

	table := &Table{
		ColumnCount: columnCount,
		Paddings:    make([]int, columnCount),
		Columns:     columns,
	}

	return table
}

// Creates a new table row with the given data and
// appends an empty row to it so it acts as a spacer
// between the header and the content.
func (t *Table) Header(data ...string) {
	t.addRow(data)

	spacer := make([]string, t.ColumnCount)
	t.addRow(spacer)
}

// Creates a new table row with the given data.
func (t *Table) Row(data ...string) {
	t.addRow(data)
}

// Displays the entire table. This will pad every string properly
// so all columns are equal as well as align everything that needs
// alignment.
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

// Adds a new row to the table and updates
// the padding based on the new content length.
// If it's longer than our longest padding, the padding
// becomes this new length.
func (t *Table) addRow(row []string) {
	for i := 0; i < t.ColumnCount; i++ {
		item := row[i]
		clean := cleanLen(item)

		// If the current element is longer than the current
		// padding, update the padding
		if clean > t.Paddings[i] {
			t.Paddings[i] = clean
		}

		// Add it to the correct column
		t.Columns[i] = append(t.Columns[i], item)
	}
}

// Pads the given string with the given padding length.
// This makes sure to account for any invisible characters
// that a string might contain and adjust the padding
// accordingly.
func formatString(str string, pad int) string {
	clean := cleanLen(str)

	width := pad + 2 // the 2 is to have an extra space on either side
	if clean != len(str) {
		width += (len(str) - clean)
	}

	// Pad the string
	final := fmt.Sprintf("%-"+strconv.Itoa(width)+"s", str)

	return final
}

// Removes all the ANSI escape characters fom the string
// as well as other werid invisible characters that a
// terminal might not print.
func cleanLen(str string) int {
	clean, _ := ansi.Strip([]byte(str))
	return len(clean)
}
