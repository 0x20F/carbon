package types

import "time"

// Container model for our own database
// specification of a container.
type Container struct {
	Id          int64     // Primary key from the database
	Uid         string    // The unique container name that we generate when starting the container
	Name        string    // The name of the service in the compose file
	ComposeFile string    // The compose file this container belongs to
	CreatedAt   time.Time // Creation time of the container
}
