package helpers

import (
	"io/ioutil"
	"os"
)

// Writes a new file to the given path, with the given name, and
// contents.
//
// It does not add any extensions to the file so it has to be
// provided in the name. Keep that in mind.
func WriteFile(path string, name string, contents []byte) (string, error) {
	// Create the file
	err := ioutil.WriteFile(path+"/"+name, contents, 0644)
	if err != nil {
		return "", err
	}

	return path + "/" + name, nil
}

// Deletes a file at the given path.
func DeleteFile(path string) error {
	return os.Remove(path)
}
