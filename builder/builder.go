package builder

import "sort"

// A Segment is the building block that all the builders
// use to compose different commands.
//
// They all define their own segments and keep track of
// the correct prioritites.
type Segment struct {
	Priority int
	Key      string
	Value    string
}

// Something that all the builders should implement,
// so that there's a generalized way that they all build
// compose themselves.
type Command interface {
	Build() string
}

// A default way of sorting all the segments based on their
// priority and making sure they get inserted properly.
//
// If the segment has a key, the value should be appended after the
// key and separated by a space.
//
// If the segment has no key, the value should be appended to the
// command normally with just a space.
func BuildCommand(command string, segments ...Segment) string {
	// Sort segments based on priority
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].Priority < segments[j].Priority
	})

	for _, segment := range segments {
		if segment.Key != "" {
			command += " " + segment.Key
		}

		if segment.Value != "" {
			command += " " + segment.Value
		}
	}

	return command
}
