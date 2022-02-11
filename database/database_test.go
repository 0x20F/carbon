package database

import (
	"co2/helpers"
	"co2/types"
	"log"
	"testing"
)

func cleanup() {
	err := helpers.DeleteFile("./carbon.db")
	if err != nil {
		log.Fatalf("Expected no error, got %s", err)
	}
}

func TestTableCreationIfNotExists(t *testing.T) {
	// Create a new database
	_, close := Get()

	defer cleanup()
	defer close()

	// Insert a new container
	InsertContainer(types.Container{})
}

func TestContainerInsert(t *testing.T) {
	// Create a new database
	_, close := Get()

	defer cleanup()
	defer close()

	// Insert a new container
	InsertContainer(types.Container{})
	InsertContainer(types.Container{})
	InsertContainer(types.Container{})

	// Get all the containers
	containers := Containers()

	// Make sure there's 3
	if len(containers) != 3 {
		t.Errorf("Expected 3 containers, got %d", len(containers))
	}
}

func TestContainerDelete(t *testing.T) {
	// Create a new database
	_, close := Get()

	defer cleanup()
	defer close()

	// Insert a new container
	InsertContainer(types.Container{Name: "test1"})
	InsertContainer(types.Container{Name: "test2"})
	InsertContainer(types.Container{Name: "test3"})

	// Get all the containers
	containers := Containers()

	// Make sure there's 3
	if len(containers) != 3 {
		t.Errorf("Expected 3 containers, got %d", len(containers))
	}

	// Delete the last container
	DeleteContainer(containers[2])

	// Get all the containers
	containers = Containers()

	// Make sure there's 2
	if len(containers) != 2 {
		t.Errorf("Expected 2 containers, got %d", len(containers))
	}
}
