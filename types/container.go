package types

import "time"

type Container struct {
	Id          int64
	Uid         string
	Name        string
	ComposeFile string
	CreatedAt   time.Time
}
