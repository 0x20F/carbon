package docker

import (
	"context"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerWrapper interface {
	RunningContainers() []dockerTypes.Container
}

type Wrapper struct{}

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
