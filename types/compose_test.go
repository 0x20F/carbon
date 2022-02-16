package types

import (
	"co2/helpers"
	"fmt"
	"testing"
)

func TestGenerateNameWorksCorrectly(t *testing.T) {
	// Build a fake compose file.
	composeFile := ComposeFile{
		Name:     "carbon.docker-compose.yml",
		Version:  "3",
		Services: make(ServiceDefinition),
	}

	// Add some services
	composeFile.Services["foo"] = make(ServiceFields)
	composeFile.Services["bar"] = make(ServiceFields)

	// Hash the names, make sure they're sorted
	hashed := helpers.Hash("bar"+"foo", 10)
	expected := fmt.Sprintf("%s.%s", hashed, composeFile.Name)

	// Check the generated name
	actual := composeFile.GenerateName()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestPathUsesGeneratedName(t *testing.T) {
	// Build a fake compose file.
	composeFile := ComposeFile{
		Name:     "carbon.docker-compose.yml",
		Version:  "3",
		Services: make(ServiceDefinition),
	}

	// Add some services
	composeFile.Services["foo"] = make(ServiceFields)
	composeFile.Services["bar"] = make(ServiceFields)

	// Hash the names, make sure they're sorted
	hashed := helpers.Hash("bar"+"foo", 10)
	name := fmt.Sprintf("%s.%s", hashed, composeFile.Name)
	expected := fmt.Sprintf("%s/%s", helpers.ComposeDir(), name)

	// Check the generated name
	actual := composeFile.Path()

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
