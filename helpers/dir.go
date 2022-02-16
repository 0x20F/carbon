package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

// Rough implementation of the user's home directory.
//
// It's a different directory depending on the platform
// you're running on so it has to be a bit more in depth.
//
// Will do some magic for Windows and will return HOME
// for all the normal OS:es
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

// Generates the path where all the carbon
// specific configurations should be stored.
//
// All generated docker compose files will be stored in
// the `~/.carbon` directory. Where that is depends on your
// platform.
//
// If the directory doesn't already exist, this will make sure
// to create it with the right permissions.
func ComposeDir() string {
	home := UserHomeDir()

	if _, err := os.Stat(home + "/.carbon"); os.IsNotExist(err) {
		os.Mkdir(home+"/.carbon", 0755)
	}

	return home + "/.carbon"
}

// Generates the path where the database file should be
// stored.
//
// Since we don't want the file to be stored wherever the binary
// is we have to store it somewhere else and the best place
// is where all the other carbon related things are. In ~/.carbon.
//
// If the directory doesn't already exist when this gets called,
// it will be created.
func DatabaseFile() string {
	home := UserHomeDir()

	if _, err := os.Stat(home + "/.carbon"); os.IsNotExist(err) {
		os.Mkdir(home+"/.carbon", 0755)
	}

	return home + "/.carbon/database.db"
}

// Turns a relative path into an absolute path.
//
// Meaning something like `./foo` will be
// turned into `/home/user/foo`.
func ExpandPath(path string) string {
	// Turn relative into absolute
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return abs
}

// Checks whether or not the provided path
// is a directory.
func IsDirectory(path string) bool {
	info, err := os.Stat(ExpandPath(path))
	if err != nil {
		panic(err)
	}

	return info.IsDir()
}
