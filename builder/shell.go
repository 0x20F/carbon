package builder

// Command builder for a `docker exec` command.
//
// Supported segments:
// - `container`: The container to be used.
// - `shell`: The shell to be used.
type DockerShellCommandBuilder struct {
	Command  string
	Segments []Segment
	Unique   map[int]Segment
}

// The container ID or name try and run the command in.
func (c *DockerShellCommandBuilder) Container(container string) *DockerShellCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      container,
		Value:    "",
	})

	return c
}

// A command to be executed within the provided container.
// This is usually used as a shell so `/bin/bash` is used by default.
//
// There's nothing that stops it from being used with other commands if
// execution on multiple containers might be needed someday.
func (c *DockerShellCommandBuilder) Shell(shell string) *DockerShellCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      shell,
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

// Creates a new instance of a Docker Exec Builder which can
// be used to dynamically build a docker exec command.
func DockerShellCommand() *DockerShellCommandBuilder {
	return &DockerShellCommandBuilder{
		Command:  "docker exec -it",
		Segments: []Segment{},
		Unique:   map[int]Segment{},
	}
}

// Interface implementation
func (c *DockerShellCommandBuilder) Build() string {
	// Get all the unique segments into an array
	segments := []Segment{}
	for _, segment := range c.Unique {
		segments = append(segments, segment)
	}

	// Merge the unique segments with the rest
	segments = append(segments, c.Segments...)

	return BuildCommand(c.Command, segments...)
}
