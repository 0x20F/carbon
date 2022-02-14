package builder

type DockerShellCommandBuilder struct {
	Command  string
	Segments []Segment
	Unique   map[int]Segment
}

func (c *DockerShellCommandBuilder) Container(container string) *DockerShellCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      container,
		Value:    "",
	})

	return c
}

func (c *DockerShellCommandBuilder) Bash() *DockerShellCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "bash",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func (c *DockerShellCommandBuilder) Sh() *DockerShellCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "sh",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func (c *DockerShellCommandBuilder) Shell(shell string) *DockerShellCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      shell,
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func DockerShellCommand() *DockerShellCommandBuilder {
	return &DockerShellCommandBuilder{
		Command:  "docker exec -it",
		Segments: []Segment{},
		Unique:   map[int]Segment{},
	}
}

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
