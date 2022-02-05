package builder

import (
	"sort"
)

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
		Priority: 1,
		Key:      "-f",
		Value:    file,
	})

	return c
}

func (c *DockerBuildCommandBuilder) Tag(tag string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 2,
		Key:      "-t",
		Value:    tag,
	})

	return c
}

func (c *DockerBuildCommandBuilder) BuildArg(arg string) *DockerBuildCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 3,
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
	var command = c.Command

	// Sort segments based on priority
	sort.Slice(c.Segments, func(i, j int) bool {
		return c.Segments[i].Priority < c.Segments[j].Priority
	})

	for _, segment := range c.Segments {
		if segment.Key != "" {
			command += " " + segment.Key
		}

		if segment.Value != "" {
			command += " " + segment.Value
		}
	}

	return command
}
