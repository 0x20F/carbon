package cmd

import (
	"co2/builder"
	"co2/carbon"
	"co2/database"
	"co2/helpers"
	"co2/runner"
	"co2/types"
	"fmt"
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
	help := "Force the start of the service. Even if it's already in the database. This will delete the old ones."
	startCmd.Flags().BoolVarP(&force, "force", "f", false, help)
}

func shouldRun(provided []string) bool {
	if force {
		return true
	}

	// Get all containers from the database
	containers := database.Containers()

	// If an of the provided containers is in the database, quit
	for _, container := range containers {
		if helpers.Contains(provided, container.Name) {
			fmt.Printf("%s is already in the database\n", container.Name)
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
		// Run the stop command
		execStop(cmd, args)
	}

	configs := carbon.Configurations("../", 2) // FIXME: Don't hardcode this
	choices := types.CarbonConfig{}

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

	// Save the compose file
	channel := make(chan bool)
	go compose.Save(channel)

	<-channel

	// Build the command
	command := builder.DockerComposeCommand().
		File(compose.Path()).
		Service(strings.Join(args, " ")).
		Background().
		Up().
		Build()

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
