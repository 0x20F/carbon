package docker

import (
	"context"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Docker API wrapper to allow for easy mocking.
type DockerWrapper interface {
	RunningContainers() []dockerTypes.Container
}

type Wrapper struct{}

// Pull the running containers directly from the docker api.
// We want speed, and that seems to be the fastest option here since
// docker itself uses this api.
func (w *Wrapper) RunningContainers() []dockerTypes.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), dockerTypes.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}
