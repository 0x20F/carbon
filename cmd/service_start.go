package cmd

import (
	"co2/builder"
	"co2/carbon"
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	force bool

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the provided services",
		Args:  cobra.MinimumNArgs(1),
		Run:   execStart,
	}
)

func init() {
	help := "Force the start of the service. This will delete the old ones before starting."
	startCmd.Flags().BoolVarP(&force, "force", "f", false, help)
}

// Checks whether or not the containers are already
// in the database. If they are, we don't want to
// start them again.
//
// If the force flag is provided, this will always return
// true./
func shouldRun(provided []string) bool {
	if force {
		return true
	}

	// Get all containers from the database
	containers := database.Containers()

	// If an of the provided containers is in the database, quit
	for _, container := range containers {
		if helpers.Contains(provided, container.Name) {
			printer.Error("ERROR", "service already running:", container.Name)
			printer.Extra(
				printer.Red,
				"Try running `co2 stop` first",
				"Use the `--force` flag to force the start",
			)
			return false
		}
	}

	return true
}

func execStart(cmd *cobra.Command, args []string) {
	if ok := shouldRun(args); !ok {
		return
	}

	// If we're forcing, we want to cleanup first
	if force {
		printer.Info(
			printer.Yellow,
			"WARN",
			"`--force` flag is set, stopping all provided services:",
			strings.Join(args, ", "),
		)
		execStop(cmd, args)
	}

	printer.Info(
		printer.Green,
		"START",
		"Starting provided services:",
		strings.Join(args, ", "),
	)
	printer.Extra(printer.Green, "Looking through the store")

	configs := carbon.Configurations("../", 2) // FIXME: Don't hardcode this
	choices := types.CarbonConfig{}

	printer.Extra(printer.Green, "Generating compose file")

	for _, service := range args {
		if _, ok := configs[service]; !ok {
			continue
		}

		found := configs[service]

		// Inject a random container name
		container := service + "-" + helpers.RandomAlphaString(10)
		found.FullContents["container_name"] = container
		found.Container = container

		choices[service] = found
	}

	compose := types.NewComposeFile()

	// Add all the services to the compose file
	for _, service := range choices {
		compose.Services[service.Name] = service.FullContents
	}

	printer.Extra(printer.Green, "Saving compose file at `"+compose.Path()+"`")

	// Save the compose file
	channel := make(chan bool)
	go compose.Save(channel)

	// Build the command
	command := builder.DockerComposeCommand().
		File(compose.Path()).
		Service(strings.Join(args, " ")).
		Background().
		Up().
		Build()

	<-channel

	printer.Extra(printer.Green, "Executing `docker compose` command on the new file\n")
	runner.Execute(command)
	go containerize(channel, compose)

	<-channel
}

// Creates container types for each of the provided services
// within the compose file, making sure that all the containers
// know which file they belong to.
//
// Also saves all the containers to the database so that all required
// information can be retrieved later if ever needed.
func containerize(channel chan bool, compose types.ComposeFile) {
	// Create container types for all services in the compose file
	containers := []types.Container{}

	for name, service := range compose.Services {
		container := types.Container{
			Name:        name,
			Uid:         service["container_name"].(string),
			ComposeFile: compose.Path(),
		}

		containers = append(containers, container)
	}

	// Create a workgroup for all the containers
	var wg sync.WaitGroup

	// Save all the containers
	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()

			// Save the container
			database.InsertContainer(container)
		}(container)
	}

	wg.Wait()
	channel <- true
}
