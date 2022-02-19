package carbon

import (
	"co2/types"
	"testing"
)

// Note that this is indented with spaces.
// Tabs will break yaml so be careful.
var customDocument = `
test:
    image: golang
    depends_on: 
        - test-db

---

test-db:
    image: lmao
    ports:
        - "8080:80"
`

func TestMovingDataBetweenDefinitions(t *testing.T) {
	// Create a new service definition
	def := types.ServiceDefinition{
		"service": {
			"unique-but-useless": "value",
		},
	}

	// Create a new carbon config
	config := types.CarbonConfig{
		"service": {},
	}

	// Move the data from the service definition into the carbon config
	k, v := move(def, config, "filename")

	// Assert that the key is the same as the service name
	if k != "service" {
		t.Errorf("Expected key to be 'service', got '%s'", k)
	}

	// Make sure the path got set to filename
	if v.Path != "filename" {
		t.Errorf("Expected path to be 'filename', got '%s'", v.Path)
	}

	// Make sure the name got set to service
	if v.Name != "service" {
		t.Errorf("Expected name to be 'service', got '%s'", v.Name)
	}

	// Make sure the unique-but-useless is present in the full contents
	if v.FullContents["unique-but-useless"] != "value" {
		t.Errorf("Expected 'unique-but-useless' to be 'value', got '%s'", v.FullContents["unique-but-useless"])
	}
}

func TestYamlParsingOfMultipleDocuments(t *testing.T) {
	// Parse the yaml into a carbon config
	config := documents([]byte(customDocument), "filename")

	// Make sure the config has the correct number of services
	if len(config) != 2 {
		t.Errorf("Expected 2 services, got %d", len(config))
	}

	// Make sure the path is filename for each of the services
	if config["test"].Path != "filename" {
		t.Errorf("Expected path to be 'filename', got '%s'", config["test"].Path)
	}

	if config["test-db"].Path != "filename" {
		t.Errorf("Expected path to be 'filename', got '%s'", config["test-db"].Path)
	}

	// Make sure the config has the correct number of fields
	if len(config["test"].FullContents) != 2 {
		t.Errorf("Expected 1 field, got %d", len(config["test"].FullContents))
	}

	// Make sure the config has the correct number of fields
	if len(config["test-db"].FullContents) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(config["test-db"].FullContents))
	}
}
