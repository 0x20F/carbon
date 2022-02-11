package cmd

import (
	"co2/docker"
	"fmt"

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

	for key, container := range containers {
		fmt.Println(key, container.Names[0], container.ID, container.Image, container.Status, container.Created)
	}
}
