package builder

// Command builder for a `docker build` command.
//
// Supported segments:
// - `path`: Path to the Dockerfile.
// - `file`: What file should be used as Dockerfile.
// - `tag`: Tag to be used for the image.
// - `build-arg`: Build argument to be passed to the image.
type DockerBuildCommandBuilder struct {
	Command  string
	Segments []Segment
}

// The path to the Dockerfile. (Think the dot in `docker build .`)
func (c *DockerBuildCommandBuilder) Path(path string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 999,
		Key:      "",
		Value:    path,
	})

	return c
}

// `-f` The file to look for in the provided Path
func (c *DockerBuildCommandBuilder) File(file string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    file,
	})

	return c
}

// `-t` The tag to be used for the image.
func (c *DockerBuildCommandBuilder) Tag(tag string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 20,
		Key:      "-t",
		Value:    tag,
	})

	return c
}

// `--build-arg` Build argument to be passed to the image.
func (c *DockerBuildCommandBuilder) BuildArg(arg string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 30,
		Key:      "--build-arg",
		Value:    arg,
	})

	return c
}

// Creates a new instance of a Docker Build Builder which can
// be used to dynamically build a docker build command.
func DockerBuildCommand() *DockerBuildCommandBuilder {
	return &DockerBuildCommandBuilder{
		Command: "docker build",
		Segments: []Segment{
			{},
		},
	}
}

// Interface implementation
func (c *DockerBuildCommandBuilder) Build() string {
	return BuildCommand(c.Command, c.Segments...)
}
