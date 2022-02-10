package builder

type DockerComposeCommandBuilder struct {
	Command  string
	Segments []Segment
	Unique   map[int]Segment
}

func (c *DockerComposeCommandBuilder) File(file string) *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    file,
	})

	return c
}

func (c *DockerComposeCommandBuilder) Background() *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 999,
		Key:      "-d",
		Value:    "",
	})

	return c
}

func (c *DockerComposeCommandBuilder) Service(service string) *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 1000,
		Key:      service,
		Value:    "",
	})

	return c
}

func (c *DockerComposeCommandBuilder) Up() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "up",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func (c *DockerComposeCommandBuilder) Down() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "down",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func (c *DockerComposeCommandBuilder) Restart() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "restart",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

func DockerComposeCommand() *DockerComposeCommandBuilder {
	return &DockerComposeCommandBuilder{
		Command: "docker compose",
		Segments: []Segment{
			{},
		},
		Unique: map[int]Segment{},
	}
}

func (c *DockerComposeCommandBuilder) Build() string {
	// Get all the unique segments into an array as well
	segments := []Segment{}
	for _, segment := range c.Unique {
		segments = append(segments, segment)
	}

	// Merge the unique segments with the actual segments
	segments = append(segments, c.Segments...)

	return BuildCommand(c.Command, segments...)
}
