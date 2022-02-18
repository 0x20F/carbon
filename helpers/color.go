package helpers

import "fmt"

// Generates a new hexadecimal color code in the format
// #RRGGBB based on the given string.
//
// Whatever the string is, the returned color will be unique
// to that string meaning that if you send the same string 50 times,
// you can expect to receive the same color 50 times.
func StringToColor(s string) string {
	hash := 0
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		hash = int(runes[i]-'0') + ((hash << 5) - hash)
	}

	color := "#"

	for i := 0; i < 3; i++ {
		value := (hash >> (i * 8)) & 0xFF
		temp := fmt.Sprintf("00%x", value)
		color += temp[len(temp)-2:]
	}

	return color
}
