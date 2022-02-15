package cmd

import (
	"co2/builder"
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the provided services",
	Args:  cobra.MinimumNArgs(1),
	Run:   execStop,
}

// Handle the stopping of the provided running containers.
// This will make sure that each of the compose files will be found
// for the provided containers, and every single one of them will be
// stopped based on their individual compose files.
//
// While that's going on it will also make sure to remove the containers
// from the database since they are technically not running anymore.
func execStop(cmd *cobra.Command, args []string) {
	printer.Info(
		printer.Green,
		"STOP",
		"Stopping provided services:",
		strings.Join(args, ", "),
	)

	containers := database.Containers()
	grouped := map[string][]types.Container{}

	// Group the containers by their compose file
	for _, container := range containers {
		if helpers.Contains(args, container.Name) {
			grouped[container.ComposeFile] = append(grouped[container.ComposeFile], container)
		}
	}

	if len(grouped) == 0 {
		printer.Extra(printer.Cyan, "None of the provided services are running", "Ignoring")
		return
	}

	var wg sync.WaitGroup

	// Stop all the containers in each group and
	// delete them from the database
	for path, composeFile := range grouped {
		command := builder.DockerComposeCommand().
			File(composeFile[0].ComposeFile).
			Stop()

		for _, container := range composeFile {
			command.Service(container.Name)

			wg.Add(1)
			go func(container types.Container) {
				defer wg.Done()
				database.DeleteContainer(container)
			}(container)
		}

		// Run the command even if the database hasn't fully
		// updated yet. It's independent.
		printer.Extra(printer.Green, "Executing stop command for compose file: "+path)
		runner.Execute(command.Build())
	}

	wg.Wait()
}
