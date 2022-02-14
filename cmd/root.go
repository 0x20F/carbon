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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(storeCmd)
	rootCmd.AddCommand(shellCmd)
}
