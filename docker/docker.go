package docker

import (
	"context"
	"sync"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var once sync.Once
var instance *impl

type impl struct {
	docker DockerWrapper
}

func wrapper() *impl {
	if instance != nil {
		return instance
	}

	return CustomWrapper(&Wrapper{})
}

func CustomWrapper(docker DockerWrapper) *impl {
	once.Do(func() {
		instance = &impl{
			docker: &Wrapper{},
		}
	})

	return instance
}

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
