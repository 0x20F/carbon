package helpers

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileCreationWithProvidedValues(t *testing.T) {
	path := "."
	name := "test_file"
	contents := []byte("test")

	path, err := WriteFile(path, name, contents)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if path != "./test_file" {
		t.Errorf("Expected %s, got %s", "./test_file", path)
	}

	// Open the file and look at the contents
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if string(fileContents) != "test" {
		t.Errorf("Expected %s, got %s", "test", string(fileContents))
	}

	file.Close()
}

func TestFileDeletion(t *testing.T) {
	// Remove the file from the previous test
	err := DeleteFile("./test_file")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
