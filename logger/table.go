package logger

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
	ColumnCount int               // How many columns the table will have
	Paddings    []int             // The widths of each column
	Alignment   map[int]Alignment // The alignment in each of the columns
	Columns     [][]string        // The data to be printed
}

func NewTable(columnCount int, alignment map[int]Alignment) *Table {
	table := &Table{
		ColumnCount: columnCount,
		Paddings:    make([]int, columnCount),
		Alignment:   alignment,
		Columns:     make([][]string, columnCount),
	}

	return table
}
