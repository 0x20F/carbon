package docker

import (
	"sync"
)

var once sync.Once
var instance *impl

// Simple implementation to make the wrapper
// accessible without generating an abundance of instances.
type impl struct {
	docker DockerWrapper
}

// Will either create a new instance of the implementation
// or return the existing one if it's already created.
func wrapper() *impl {
	if instance != nil {
		return instance
	}

	return CustomWrapper(&Wrapper{})
}

// Creates a custom wrapper with the given docker wrapper
// implementation.
//
// Note that this is public so that it's usable in tests as well.
// That might not be necessary but it's open if ever needed.
func CustomWrapper(docker DockerWrapper) *impl {
	once.Do(func() {
		instance = &impl{
			docker: docker,
		}
	})

	return instance
}
