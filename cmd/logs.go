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

// Set up all the required flags.
func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "follow the logs")
}

// Filters all the available containers based on the
// provided list of IDs or service names, and then generates
// specific docker commands to be run simultaneously.
//
// If none of the provided IDs or service names match any
// of the existing containers, we don't do anything. Just inform
// the user.
func execLogs(cmd *cobra.Command, args []string) {
	matches := filterContainers(args)
	commands := generateCommands(matches, follow)

	if !shouldRunLogsCommand(commands) {
		printer.Error("ERROR", "no containers found", strings.Join(args, ", "))
		return
	}

	<-runner.Execute(commands...)
}

// Various gates and checks to make sure the
// generated commands are fit for running.
func shouldRunLogsCommand(commands []types.Command) bool {
	return len(commands) != 0
}

// Generates a list of commands to run based on the
// arguments that the user has provided.
//
// This will build docker logs commands for each of the
// containers it has been provided.
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

// Filters the list of service names and IDs provided
// by the user and tries to match them to a running container
// instance.
//
// This will look both at all the running containers and the carbon-specifit ones
// stored in the database in order to find any non-carbon containers that are
// running, and also look.
//
// It will first look at all the UIDs of the running containers since those
// don't have any specific carbon service names. If none of those match
// it will try to match the strings with the service names of the
// containers stored in the database.
//
// If a container UID is provided, as long as that container is running, it won't
// matter if it's carbon or not. As soon as a service name is provided, the service
// has to be a carbon service.
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

	// Only check for the rest if we didn't find all the UIDs
	if len(matches) == len(choices) {
		return matches
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
