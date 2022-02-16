package cmd

import (
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Removes a store with the given ID",
		Args:  cobra.MinimumNArgs(1),
		Run:   execRemove,
	}
)

func init() {}

// Looks through all the registered store UIDs
// and removes all the ones that are registered
// with the UIDs provided by the user.
func execRemove(cmd *cobra.Command, args []string) {
	stores := database.Stores()

	printer.Info(
		printer.Green,
		"REMOVE",
		fmt.Sprintf("Removing %d stores", len(args)),
		strings.Join(args, ", "),
	)

	for _, store := range stores {
		if helpers.Contains(args, store.Uid) {
			database.DeleteStore(store)
			printer.Extra(printer.Green, "Removed store: "+store.Uid)
		}
	}
}
