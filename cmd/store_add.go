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

// Adds all the required flags
func init() {
	addCmd.Flags().StringVarP(&store, "store", "s", "", "The path to the store")
	addCmd.Flags().StringVarP(&id, "id", "i", "", "The id of the store. If left empty, it will be generated")
	addCmd.Flags().StringVarP(&env, "env", "e", "", "The environment file to use for this store. Should a path to the .env file.")
}

// Registers a new carbon store
//
// This won't register anything if the store is undefined, or the
// provided store path is not a directory.
//
// This will also not register anytthing if the provided
// environment file path is a directory and not a file.
//
// It will also generate a unique identifier based on the
// store path so that it can be accessed easily later on
// when deleting, as long as the user hasn't provided their
// own identifier.
//
// If a store with an identical identifier exists already, it will
// be deleted first. Dems the rules... No duplicates.
func execAdd(cmd *cobra.Command, args []string) {
	if store == "" || !helpers.IsDirectory(store) {
		printer.Error("ERROR", "No store directory", "")
		printer.Extra(printer.Red, "You must provide a store directory with `--store`")

		if !helpers.IsDirectory(store) {
			printer.Extra(printer.Red, "The provided directory isn't actually a directory!")
		}

		return
	}

	if env != "" && helpers.IsDirectory(env) {
		printer.Error("ERROR", "Environment file is wrong", "")
		printer.Extra(printer.Red, "The provided file isn't actually a file!")

		return
	}

	if id == "" {
		printer.Info(printer.Cyan, "INFO", "No id provided for store:", store)
		printer.Extra(printer.Cyan, "Generating a random id for you")

		id = helpers.Hash(store, 4)
	}

	store = helpers.ExpandPath(store)
	environment := ""
	printer.Info(printer.Green, "ADD", "Adding store", store)

	if env != "" {
		environment = helpers.ExpandPath(env)
	}

	store := types.Store{
		Uid:  id,
		Path: store,
		Env:  environment,
	}

	// Delete if it exists before adding again
	database.DeleteStore(store)
	database.AddStore(store)

	printer.Extra(
		printer.Green,
		"The id for the new store is: "+id,
		"Use `co2 show --stores` to see all id's",
		"Verify all your services are found with `co2 show -c`",
	)

	if env == "" {
		printer.Extra(printer.Cyan, "No environment file provided. No environment variables will be available for this store.")
	}
}
