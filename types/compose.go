package types

import (
	"co2/helpers"
	"fmt"
)

type ComposeFile struct {
	Name     string            `yaml:"-"`
	Version  string            `yaml:"version"`
	Services ServiceDefinition `yaml:"services"`
}

func NewComposeFile() ComposeFile {
	name := randomComposeName()

	return ComposeFile{
		Name:     name,
		Version:  "3",
		Services: make(ServiceDefinition),
	}
}

func randomComposeName() string {
	return fmt.Sprintf("%s.docker-compose.yml", helpers.RandomAlphaString(10))
}
