package docker

import "co2/types"

func RunningContainers() []string {
	return []string{}
}

func NewComposeFile() types.ComposeFile {
	return types.ComposeFile{
		Version:  "3",
		Services: make(types.ServiceDefinition),
	}
}
