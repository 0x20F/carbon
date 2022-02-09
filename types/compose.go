package types

import (
	"co2/helpers"

	"gopkg.in/yaml.v2"
)

type ComposeFile struct {
	Name     string            `yaml:"-"`
	Version  string            `yaml:"version"`
	Services ServiceDefinition `yaml:"services"`
}

func (c *ComposeFile) Path() string {
	return helpers.ComposeDir() + "/" + c.Name
}

func (c *ComposeFile) Save(finished chan bool) {
	// Turn the ComposeFile into a string
	contents, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}

	// Save it to the file
	savePath := helpers.ComposeDir()
	_, err = helpers.WriteFile(savePath, c.Name, contents)
	if err != nil {
		panic(err)
	}

	finished <- true
}
