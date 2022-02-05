package builder

type DockerBuildCommandBuilder struct {
	Command  string
	Segments []Segment
}

func (c *DockerBuildCommandBuilder) Path(path string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 999,
		Key:      "",
		Value:    path,
	})

	return c
}

func (c *DockerBuildCommandBuilder) File(file string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    file,
	})

	return c
}

func (c *DockerBuildCommandBuilder) Tag(tag string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 20,
		Key:      "-t",
		Value:    tag,
	})

	return c
}

func (c *DockerBuildCommandBuilder) BuildArg(arg string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 30,
		Key:      "--build-arg",
		Value:    arg,
	})

	return c
}

func DockerBuildCommand() *DockerBuildCommandBuilder {
	return &DockerBuildCommandBuilder{
		Command: "docker build",
		Segments: []Segment{
			{},
		},
	}
}

func (c *DockerBuildCommandBuilder) Build() string {
	return BuildCommand(c.Command, c.Segments...)
}
