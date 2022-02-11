package types

import (
	"co2/helpers"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type ComposeFile struct {
	Name          string            `yaml:"-"`        // The name of the compose file without the unique id
	Version       string            `yaml:"version"`  // The version of the compose file
	Services      ServiceDefinition `yaml:"services"` // A map of all services this compose file contains
	GeneratedName string            `yaml:"-"`        // The name of the compose file with the unique id prepended
}

func NewComposeFile() ComposeFile {
	return ComposeFile{
		Name:     "carbon.docker-compose.yml",
		Version:  "3",
		Services: make(ServiceDefinition),
	}
}

// Generate a unique identifier for the compose file
// based on the name of the services that the compose file
// contains.
//
// If the same services get booted up multiple times, we want the
// compose files to be reused not recreated with a whole new name.
// This makes us have to worry a lot less about cleaning up the
// orphaned compose files since there won't be that many.
func (c *ComposeFile) GenerateName() string {
	if c.GeneratedName != "" {
		return c.GeneratedName
	}

	// Join all the service names together
	var names []string
	for name := range c.Services {
		names = append(names, name)
	}

	// Sort them so they're always in the same order
	sort.Strings(names)

	// Build the final name
	hash := helpers.Hash(strings.Join(names, ""), 10)
	name := fmt.Sprintf("%s.%s", hash, c.Name)
	c.GeneratedName = name

	return c.GeneratedName
}

// Concat method for adding the default carbon compose
// configuration directory and the unique compose file name
// based on all the services together.
func (c *ComposeFile) Path() string {
	return helpers.ComposeDir() + "/" + c.GenerateName()
}

// Convert the file into yaml and save it to its designated
// path.
//
// This accepts a channel since it might be something that needs
// to happen but we don't really care about the result as long as
// it works so we can turn it into a goroutine instead.
func (c *ComposeFile) Save(finished chan bool) {
	contents, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}

	savePath := helpers.ComposeDir()
	// Save it to the file
	_, err = helpers.WriteFile(savePath, c.GenerateName(), contents)
	if err != nil {
		panic(err)
	}

	finished <- true
}
