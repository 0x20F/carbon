package docker

import (
	"co2/types"
	"fmt"
	"strings"
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
func RunningContainers() []types.Container {
	cli := wrapper().docker
	containers := cli.RunningContainers()

	var parsed = []types.Container{}

	for _, container := range containers {
		ports := []string{}

		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", port.PublicPort, port.Type))
		}

		c := types.Container{
			Name:      container.Names[0][1:], // Remove the leading /
			Image:     container.Image,
			Ports:     strings.Join(ports, ", "),
			Status:    container.Status,
			DockerUid: container.ID,
		}
		c.Hash()

		parsed = append(parsed, c)
	}

	return parsed
}
