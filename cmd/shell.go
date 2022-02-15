package cmd

import (
	"co2/builder"
	"co2/docker"
	"co2/printer"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	sh     bool
	custom string

	shellCmd = &cobra.Command{
		Use:   "shell",
		Short: "Composes a command for getting a shell within a specific container by id",
		Args:  cobra.MinimumNArgs(1),
		Run:   execShell,
	}
)

// Adds the required flags
func init() {
	shellCmd.Flags().BoolVarP(&sh, "sh", "s", false, "Use sh")
	shellCmd.Flags().StringVarP(&custom, "custom", "c", "", "Use custom shell")
}

// Executes the command and makes sure that the
// correct shell gets added to the built docker command.
//
// This will assume that /bin/bash is what you want most of
// the time so that's the default shell.
//
// It only looks through the running containers, and is not
// dependent of any carbon features since it runs on the
// docker api.
//
// Running the command and returning an interactive shell is too
// complex for now so the only thing this will do, if it manages
// to compile the command is to print it to standard output so it
// can be piped into something else.
func execShell(cmd *cobra.Command, args []string) {
	id := args[0]
	shell := "/bin/bash"

	if sh {
		shell = "/bin/sh"
	}

	if custom != "" {
		shell = custom
	}

	// Get all containers
	containers := docker.RunningContainers()

	for uid, container := range containers {
		if uid != id {
			continue
		}

		cmd := builder.DockerShellCommand().
			Container(container.Names[0]).
			Shell(shell).
			Build()

		// Output the command
		fmt.Println(cmd)
		return
	}

	// If we get here, we didn't find the container
	printer.Error("ERROR", "container not found:", id)
}
