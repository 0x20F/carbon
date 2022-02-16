package cmd

import (
	"github.com/spf13/cobra"
)

var (
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Does things with the services you provide",
		Args:  cobra.MinimumNArgs(1),
		Run:   execService,
	}
)

// Adds the required subcommands since
// this command does nothing on its own.
func init() {
	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(stopCmd)
}

func execService(cmd *cobra.Command, args []string) {}
