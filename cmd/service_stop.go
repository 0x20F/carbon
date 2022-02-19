package cmd

import (
	"co2/builder"
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"strings"

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

	groups := groupByComposeFile(args...)

	if len(groups) == 0 {
		printer.Extra(printer.Cyan, "None of the provided services are running", "Ignoring")
		return
	}

	stopContainers(groups)
}

// Groups all the carbon service IDs or names that the
// user has provided by their respective compose files.
// Returns a map of compose file paths to a list of containers
// that should be stopped in that compose file.
func groupByComposeFile(choices ...string) map[string][]types.Container {
	containers := database.Containers()
	groups := make(map[string][]types.Container)

	for _, container := range containers {
		if helpers.Contains(choices, container.ServiceName) ||
			helpers.Contains(choices, container.Uid) {

			groups[container.ComposeFile] = append(groups[container.ComposeFile], container)
		}
	}

	return groups
}

// Builds a new docker compose stop command for each provided
// compose file container group and then runs them all in parallel after
// deleting all the containers from the database.
func stopContainers(groups map[string][]types.Container) {
	commands := []types.Command{}

	for _, composeFile := range groups {
		command := builder.DockerComposeCommand().
			File(composeFile[0].ComposeFile).
			Stop()

		for _, container := range composeFile {
			command.Service(container.ServiceName)
			database.DeleteContainer(container)
		}

		commands = append(commands, types.Command{
			Text: command.Build(),
		})
	}

	printer.Extra(printer.Green, "Executing stop commands")
	runner.Execute(commands...)
}
