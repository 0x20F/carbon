package docker

import (
	"co2/helpers"

	dockerTypes "github.com/docker/docker/api/types"
)

func RunningContainers() map[string]dockerTypes.Container {
	cli := wrapper().docker
	containers := cli.RunningContainers()

	var parsed = map[string]dockerTypes.Container{}

	for _, container := range containers {
		key := helpers.Hash(container.Image+container.Names[0], 4)

		parsed[key] = container
	}

	return parsed
}
