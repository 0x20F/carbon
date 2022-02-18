package types

import (
	"co2/helpers"
	"time"
)

// Container model for our own database
// specification of a container.
type Container struct {
	Id          int64    // Primary key from the database
	DockerUid   string   // The ID that docker gives to each container
	Uid         string   // A unique hash for the container based on the image and name
	Name        string   // The unique container name that we generate when starting the container
	Image       string   // The image of the container
	ServiceName string   // The name of the service in the compose file
	ComposeFile string   // The compose file this container belongs to
	Ports       []string // All exposed ports
	Status      string
	CreatedAt   time.Time // Creation time of the container
}

func (c *Container) Hash() {
	c.Uid = helpers.Hash(c.Image+c.Name, 4)
}
