package builder

type DockerLogsCommandBuilder struct {
	Command  string
	Segments []Segment
}

func (c *DockerLogsCommandBuilder) Follow() *DockerLogsCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    "",
	})

	return c
}

func (c *DockerLogsCommandBuilder) Container(name string) *DockerLogsCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 100,
		Key:      name,
		Value:    "",
	})

	return c
}

func DockerLogsCommand() *DockerLogsCommandBuilder {
	return &DockerLogsCommandBuilder{
		Command:  "docker logs",
		Segments: []Segment{},
	}
}

func (c *DockerLogsCommandBuilder) Build() string {
	return BuildCommand(c.Command, c.Segments...)
}
