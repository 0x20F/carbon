package docker

import (
	"sync"
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
