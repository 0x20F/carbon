package logger

import "fmt"

func Green(title string, info string, highlight string) {
	fmt.Printf("\x1b[32m%s\x1b[0m: %s\n", title, info)
}
