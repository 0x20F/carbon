package builder

import "sort"

type Segment struct {
	Priority int
	Key      string
	Value    string
}

type Command interface {
	Build() string
}

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
