package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")

		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		return home
	}

	return os.Getenv("HOME")
}

func ComposeDir() string {
	home := UserHomeDir()

	// Check if the user has a .carbon directory in their home
	if _, err := os.Stat(home + "/.carbon"); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		os.Mkdir(home+"/.carbon", 0755)
	}

	return home + "/.carbon"
}

func ExpandPath(path string) string {
	// Turn relative into absolute
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return abs
}
