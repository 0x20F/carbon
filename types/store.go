package types

import "time"

// Store model for our own database
// specification of a store.
type Store struct {
	Id        int64     // Database key
	Uid       string    // Unique identifier for the store
	Path      string    // The path to the store
	Env       string    // The environment file linked to this store
	CreatedAt time.Time // The time the store was created at
}
