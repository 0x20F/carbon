package cmd

import (
	"co2/database"
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
	showCmd.Flags().BoolVarP(&running, "running", "r", false, "show all currently running carbon containers")
}

func execShow(cmd *cobra.Command, args []string) {
	if !running {
		fmt.Println("Nothing to show...")
		return
	}

	// Get all containers
	containers := database.Containers()

	// Print them all
	for _, container := range containers {
		fmt.Println(container)
	}
}
