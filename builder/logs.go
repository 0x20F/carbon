package builder

// Command builder for a `docker logs` command.
//
// Supported segments:
// - `follow`: Follow the logs.
// - `container`: The container to get the logs from.
type DockerLogsCommandBuilder struct {
	Command  string
	Segments []Segment
}

// `-f` Follow the logs.
func (c *DockerLogsCommandBuilder) Follow() *DockerLogsCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    "",
	})

	return c
}

// The container to get the logs from.
func (c *DockerLogsCommandBuilder) Container(name string) *DockerLogsCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 100,
		Key:      name,
		Value:    "",
	})

	return c
}

// Creates a new instance of a Docker Logs Builder which can
// be used to dynamically build a docker logs command.
func DockerLogsCommand() *DockerLogsCommandBuilder {
	return &DockerLogsCommandBuilder{
		Command:  "docker logs",
		Segments: []Segment{},
	}
}

// Interface implementation
func (c *DockerLogsCommandBuilder) Build() string {
	return BuildCommand(c.Command, c.Segments...)
}
