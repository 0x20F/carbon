package types

import (
	dockerTypes "github.com/docker/docker/api/types"
)

type SortableMapItem struct {
	Uid       string
	Container dockerTypes.Container
}

type SortableMap []SortableMapItem

type Command struct {
	Text string
	Name string
}
