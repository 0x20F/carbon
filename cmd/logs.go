package cmd

import (
	"co2/builder"
	"co2/database"
	"co2/docker"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"strings"

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
	matches := filterContainers(args)
	commands := generateCommands(matches, follow)

	if !shouldRunLogsCommand(commands) {
		printer.Error("ERROR", "no containers found", strings.Join(args, ", "))
		return
	}

	// Execute all the commands
	<-runner.Execute(commands...)
}

func shouldRunLogsCommand(commands []types.Command) bool {
	return len(commands) != 0
}

func generateCommands(matches []types.Container, follow bool) []types.Command {
	var commands = []types.Command{}

	for _, match := range matches {
		command := builder.DockerLogsCommand().
			Container(match.Name)

		if follow {
			command.Follow()
		}

		commands = append(commands, types.Command{
			Text: command.Build(),
			Name: match.Name,
		})
	}

	return commands
}

func filterContainers(choices []string) []types.Container {
	containers := docker.RunningContainers()
	saved := database.Containers()

	var matches = []types.Container{}

	// Check for Uids
	for _, structure := range containers {
		if helpers.Contains(choices, structure.Uid) {
			matches = append(matches, structure)
			continue
		}
	}

	// Check for service names
	for _, container := range saved {
		if !helpers.Contains(choices, container.ServiceName) {
			continue
		}

		matched := false

		// Check if we've already matched on the Uid in the above loop
		for _, match := range matches {
			if match.Uid == container.Uid {
				matched = true
				break
			}
		}

		if matched {
			continue
		}

		matches = append(matches, container)
	}

	return matches
}
