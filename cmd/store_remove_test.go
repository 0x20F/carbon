package cmd

import (
	"co2/database"
	"co2/types"
	"testing"
)

func TestRemoveByStoreUid(t *testing.T) {
	beforeCmdTest()

	// Clean up the database first
	for _, store := range database.Stores() {
		database.DeleteStore(store)
	}

	// Add a bunch of stores to the database
	stores := []types.Store{
		{
			Uid:  "uid1",
			Path: "path1",
		},
		{
			Uid:  "uid2",
			Path: "path2",
		},
		{
			Uid:  "uid3",
			Path: "path3",
		},
		{
			Uid:  "uid4",
			Path: "path4",
		},
	}

	for _, store := range stores {
		database.AddStore(store)
	}

	// Remove 2
	removeByUid("uid2", "uid1")

	// Make sure there's only one store in the database
	if len(database.Stores()) != 2 {
		t.Error("removeByStoreUid should remove the store with the matching uid. Expected 2 stores, got", len(database.Stores()))
	}

	afterCmdTest()
}
