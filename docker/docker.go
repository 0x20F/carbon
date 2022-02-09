package docker

import (
	"co2/helpers"
	"co2/types"
	"fmt"
)

func RunningContainers() []string {
	return []string{}
}

func NewComposeFile() types.ComposeFile {
	name := randomComposeName()

	fmt.Println("New name is:", name)

	return types.ComposeFile{
		Name:     name,
		Version:  "3",
		Services: make(types.ServiceDefinition),
	}
}

func randomComposeName() string {
	return fmt.Sprintf("%s.docker-compose.yml", helpers.RandomAlphaString(10))
}
