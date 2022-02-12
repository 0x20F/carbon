package cmd

import (
	"co2/docker"
	"co2/logger"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	running bool

	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows different kinds of information",
		Run:   execShow,
	}
)

func init() {
	showCmd.Flags().BoolVarP(&running, "running", "r", false, "show all currently running containers")
}

func execShow(cmd *cobra.Command, args []string) {
	if running {
		showRunning()
	}
}

func showRunning() {
	containers := docker.RunningContainers()
	table := logger.NewTable(7)

	fmt.Println()
	logger.Info("RUN", "total running containers:", fmt.Sprint(len(containers)))

	table.Header(
		"KEY",
		"NAME",
		"ID",
		"IMAGE",
		"PORTS",
		"CREATED",
		"STATUS",
	)

	for key, container := range containers {
		// Turn the array of ports into a string
		ports := []string{}

		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", port.PublicPort, port.Type))
		}

		table.AddRow(
			key,
			container.Names[0],
			container.ID[:7],
			container.Image,
			strings.Join(ports, ", "),
			fmt.Sprint(container.Created),
			container.Status,
		)
	}

	table.Display()
}
