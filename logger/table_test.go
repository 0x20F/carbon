package logger

import "testing"

func TestTableApi(t *testing.T) {
	table := NewTable(3)

	table.Header("Header 1", "Header 2", "Header 3")

	table.AddRow("Row 1", "Row 2", "Row 3")
	table.AddRow("Row 4", "Row 5", "Row 6")
	table.AddRow("Row 7", "A bit longer tho to test padding", "Row 9")

	table.Display()
}

func TestFormatString(t *testing.T) {
	str := "1234" // len -> 4
	pad := 10     // len -> 10
	// final len -> 4 + (10 - 4) + 2

	formatted := formatString(str, pad)

	if len(formatted) != 12 {
		t.Error("Expected len", 12, "got", len(formatted))
	}
}
