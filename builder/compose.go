package builder

// Command builder for a `docker compose` command.
//
// Supported segments:
// - `file`: What file should be used as the docker-compose.yml
// - `env-file`: What file(s) should be used as the .env file(s)
// - `background`: Should the command be run in the background
// - `service`: What service(s) should be started
// - `up`: Start the services
// - `down`: Take down the entire compose file
// - `stop`: Stop the services
// - `restart`: Restart the services
type DockerComposeCommandBuilder struct {
	Command  string
	Segments []Segment
	Unique   map[int]Segment
}

// `-f` what file should be used as the docker-compose.yml
func (c *DockerComposeCommandBuilder) File(file string) *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "-f",
		Value:    file,
	})

	return c
}

// `--env-file` what file(s) should be used as the .env file(s)
func (c *DockerComposeCommandBuilder) EnvFile(file string) *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 10,
		Key:      "--env-file",
		Value:    file,
	})

	return c
}

// `-d` should the command be run in the background
func (c *DockerComposeCommandBuilder) Background() *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 1000,
		Key:      "-d",
		Value:    "",
	})

	return c
}

// A list of space separated service names that should run
func (c *DockerComposeCommandBuilder) Service(service string) *DockerComposeCommandBuilder {
	c.Segments = append(c.Segments, Segment{
		Priority: 1001,
		Key:      service,
		Value:    "",
	})

	return c
}

// `up` for starting the service
// This is unique and competes with `down`, `stop`, and `restart`
func (c *DockerComposeCommandBuilder) Up() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "up",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

// `down` for taking down the entire compose file
// This is unique and competes with `up`, `stop`, and `restart`
func (c *DockerComposeCommandBuilder) Down() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "down",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

// `stop` for stopping one or more services
// This is unique and competes with `up`, `down`, and `restart`
func (c *DockerComposeCommandBuilder) Stop() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "stop",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

// `restart` for restarting one or more services
// This is unique and competes with `up`, `down`, and `stop`
func (c *DockerComposeCommandBuilder) Restart() *DockerComposeCommandBuilder {
	segment := Segment{
		Priority: 999,
		Key:      "restart",
		Value:    "",
	}

	c.Unique[segment.Priority] = segment

	return c
}

// Creates a new instance of a Docker Compose Builder which can
// be used to dynamically build a docker-compose command.
func DockerComposeCommand() *DockerComposeCommandBuilder {
	return &DockerComposeCommandBuilder{
		Command: "docker compose",
		Segments: []Segment{
			{},
		},
		Unique: map[int]Segment{},
	}
}

// Interface implementation
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
