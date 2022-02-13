package types

import "time"

type Store struct {
	Id        int64
	Uid       string
	Path      string
	CreatedAt time.Time
}
