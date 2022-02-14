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

func init() {
	shellCmd.Flags().BoolVarP(&sh, "sh", "s", false, "Use sh")

	shellCmd.Flags().StringVarP(&custom, "custom", "c", "", "Use custom shell")
}

func execShell(cmd *cobra.Command, args []string) {
	id := args[0]
	shell := "bash" // Bash is the default shell

	if sh {
		shell = "sh"
	}

	if custom != "" {
		shell = custom
	}

	// Get all containers
	containers := docker.RunningContainers()

	// Find the container with the provided id
	for uid, container := range containers {
		if uid != id {
			continue
		}

		// Build the docker command
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
