package cmd

import (
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/types"

	"github.com/spf13/cobra"
)

var (
	store string
	id    string
	env   string

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new item to the store",
		Run:   execAdd,
	}
)

func init() {
	addCmd.Flags().StringVarP(&store, "store", "s", "", "The path to the store")
	addCmd.Flags().StringVarP(&id, "id", "i", "", "The id of the store. If left empty, it will be generated")
	addCmd.Flags().StringVarP(&env, "env", "e", "", "The environment file to use for this store. Should a path to the .env file.")
}

func execAdd(cmd *cobra.Command, args []string) {
	if store == "" || !helpers.IsDirectory(store) {
		printer.Error("ERROR", "No store directory", "")
		printer.Extra(printer.Red, "You must provide a store directory with `--store`")

		if !helpers.IsDirectory(store) {
			printer.Extra(printer.Red, "The provided directory isn't actually a directory!")
		}

		return
	}

	if env == "" || helpers.IsDirectory(env) {
		printer.Error("ERROR", "No environment file", "")
		printer.Extra(printer.Red, "You must provide an environment file with `--env`")

		if helpers.IsDirectory(env) {
			printer.Extra(printer.Red, "The provided file isn't actually a file!")
		}

		return
	}

	if id == "" {
		printer.Info(printer.Cyan, "INFO", "No id provided for store:", store)
		printer.Extra(printer.Cyan, "Generating a random id for you")

		id = helpers.Hash(store, 4)
	}

	store = helpers.ExpandPath(store)
	printer.Info(printer.Green, "ADD", "Adding store", store)

	store := types.Store{
		Uid:  id,
		Path: store,
		Env:  env,
	}

	// Try deleting the stores with the same ID if they exist
	// before adding.
	database.DeleteStore(store)
	database.AddStore(store)

	printer.Extra(
		printer.Green,
		"The id for the new store is: "+id,
		"Use `co2 show --stores` to see all id's",
		"Verify all your services are found with `co2 show -c`",
	)
}
