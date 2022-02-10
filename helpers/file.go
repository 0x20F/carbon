package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

func WriteFile(path string, name string, contents []byte) (string, error) {
	// Print the path
	fmt.Printf("Writing %s to %s\n", name, path)

	// Create the file
	err := ioutil.WriteFile(path+"/"+name, contents, 0644)
	if err != nil {
		return "", err
	}

	return path + "/" + name, nil
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
