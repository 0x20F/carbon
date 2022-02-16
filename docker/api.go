package docker

import (
	"co2/helpers"

	dockerTypes "github.com/docker/docker/api/types"
)

// Gets all the containers that are currently running on the machine.
//
// This makes use of the golang implementation of the docker api, which
// I think is the original, so it's pretty fast. Even compared to the shell
// commands that are used to get the same information.
//
// Before all the containers are returned, this will also make sure to
// inject a unique identifier for each container. This is unique for each container
// but it's also static. Meaning that as long as the container has the same name
// as its always had and the same image, the resulting unique ID will always
// be the same.
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
