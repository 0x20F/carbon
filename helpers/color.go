package helpers

import "fmt"

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
