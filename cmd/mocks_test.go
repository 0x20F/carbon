package cmd

import (
	"co2/docker"
	"co2/helpers"

	"github.com/4khara/replica"
	dockerTypes "github.com/docker/docker/api/types"
)

type MockWrapperCmd struct{}

func (w *MockWrapperCmd) RunningContainers() []dockerTypes.Container {
	_, rv := replica.MockFn()

	if rv != nil {
		var containers []dockerTypes.Container

		if rv[0] != nil {
			containers = rv[0].([]dockerTypes.Container)
		}

		return containers
	}

	return []dockerTypes.Container{
		{
			ID:    helpers.Hash("ayeeeeee lmaoooooo", 30),
			Image: "image1",
			Names: []string{"/docker-container1"},
		},
		{
			ID:    helpers.Hash("ayeeeeee lmat1r3t13ooo", 30),
			Image: "image2",
			Names: []string{"/docker-container2"},
		},
		{
			ID:    helpers.Hash("ayeeee2po3ewmföqaee lmaoooooo", 30),
			Image: "image3",
			Names: []string{"/docker-container3"},
		},
	}
}

func beforeCmdTest() {
	docker.CustomWrapper(&MockWrapperCmd{})
	replica.Mocks.Clear()
}
