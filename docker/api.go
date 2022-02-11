package docker

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunningContainers() map[string]dockerTypes.Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), dockerTypes.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var parsed = map[string]dockerTypes.Container{}

	for _, container := range containers {
		key := containerHash(&container)

		parsed[key] = container
	}

	return parsed
}

// Generate a unique hash for a container based on the container
// image and the container name. Every time the container is started,
// as long as it has the same name and image, the identifier should
// always be the same.
//
// The reason we're running FNV here is for speed, more uniqueness might
// happen with something like md5 but that was once created for cryptographic
// purposes which means the speed was never of focus when developing it.
func containerHash(container *dockerTypes.Container) string {
	h := fnv.New32a()
	io.WriteString(h, container.Image+container.Names[0])

	// Format and pad with zeros if length isn't 4
	return fmt.Sprintf("%04x", h.Sum32())[:4]
}
