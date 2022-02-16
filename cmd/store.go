package cmd

import "github.com/spf13/cobra"

var (
	storeCmd = &cobra.Command{
		Use:   "store",
		Short: "Manages the stored directories in which carbon looks",
		Run:   execStore,
	}
)

// Registers all subcommands
func init() {
	storeCmd.AddCommand(addCmd)
	storeCmd.AddCommand(removeCmd)
}

func execStore(cmd *cobra.Command, args []string) {}
