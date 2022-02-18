package types

// Wrapper around our specific way of handling commands.
// Each command needs to have a label so that the output
// is formatted nicely.
type Command struct {
	Text  string
	Label string
}
