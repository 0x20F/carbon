package cmd

import (
	"co2/database"
	"co2/helpers"
	"testing"
)

func TestShouldAddStoreOnlyStoresDirectories(t *testing.T) {
	// Make sure it returns true for a directory
	if !shouldAddStore(helpers.UserHomeDir()) {
		t.Error("shouldAddStore should return true for a directory")
	}

	// Make sure it returns false for a file
	if shouldAddStore(helpers.DatabaseFile()) {
		t.Error("shouldAddStore should return false for a file")
	}

	// Make sure it returns false for an empty string
	if shouldAddStore("") {
		t.Error("shouldAddStore should return false for an empty string")
	}
}

func TestShouldAddEnvOnlyStoresFiles(t *testing.T) {
	file := helpers.ComposeDir() + "/test"

	// Create a file beforehand
	helpers.WriteFile(helpers.ComposeDir(), "test", []byte(""))

	// Make sure it returns true for a file
	if !shouldAddEnv(file) {
		t.Error("shouldAddEnv should return true for a file")
	}

	// Make sure it returns false for a directory
	if shouldAddEnv(helpers.UserHomeDir()) {
		t.Error("shouldAddEnv should return false for a directory")
	}

	// Remove the file
	helpers.DeleteFile(file)
}

func TestValidateIdHashesStoreWhenNoneProvided(t *testing.T) {
	store := "test"
	hash := helpers.Hash(store, 4)

	// Make sure the hash is the same
	if validateId("", store) != hash {
		t.Error("validateId should hash the store when no id is provided")
	}
}

func TestValidateIdReturnsTheIdWhenProvided(t *testing.T) {
	store := "test"
	id := "id"

	// Make sure the id is the same
	if validateId(id, store) != id {
		t.Error("validateId should return the id when provided")
	}
}

func TestAddStoreDoesNotDuplicateStoresWithTheSameUid(t *testing.T) {
	beforeCmdTest()

	// Add a bunch of identical stores
	addStore("", store, "")
	addStore("", store, "")
	addStore("", store, "")
	addStore("", store, "")

	// Make sure there's only one store in the database
	if len(database.Stores()) != 1 {
		t.Error("addStore should not duplicate stores with the same path")
	}
}
