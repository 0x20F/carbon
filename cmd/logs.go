package cmd

import (
	"co2/builder"
	"co2/database"
	"co2/docker"
	"co2/helpers"
	"co2/runner"
	"co2/types"

	"github.com/spf13/cobra"
)

var (
	follow bool

	logsCmd = &cobra.Command{
		Use:   "logs",
		Short: "Shows the logs of the provided services",
		Args:  cobra.MinimumNArgs(1),
		Run:   execLogs,
	}
)

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow the logs")
}

func execLogs(cmd *cobra.Command, args []string) {
	// Get all the running containers
	containers := docker.RunningContainers()
	saved := database.Containers()

	// Find all the container IDs that the user cares about
	var matches = []types.Container{}

	for _, structure := range containers {
		if helpers.Contains(args, structure.Uid) {
			matches = append(matches, structure)
			continue
		}

		// If we didn't match the ID check for the actual service name
		for _, container := range saved {
			if helpers.Contains(args, container.ServiceName) {
				matches = append(matches, structure)
				break
			}
		}
	}

	// For each match we found, build the docker command
	var commands = []types.Command{}

	for _, structure := range matches {
		command := builder.DockerLogsCommand().
			Container(structure.Name)

		if follow {
			command.Follow()
		}

		commands = append(commands, types.Command{
			Text: command.Build(),
			Name: structure.Name,
		})
	}

	// Execute all the commands
	<-runner.Execute(commands...)
}
