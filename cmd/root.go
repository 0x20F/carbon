package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "carbon",
		Short: "Mess around with containers!!",
		Long:  "Flip, Twist, and turn all your containers!!!",
	}
)

// Starts the initial command
// which in turn will start all the other commands.
func Execute() error {
	return rootCmd.Execute()
}

// Registers all subcommands
func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)

	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(storeCmd)
	rootCmd.AddCommand(shellCmd)
	rootCmd.AddCommand(logsCmd)
}
