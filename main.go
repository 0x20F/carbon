package main

import (
	"co2/cmd"
	"math/rand"
	"time"
)

func main() {
	// Seed random for everything that needs it
	rand.Seed(time.Now().UnixNano())

	// Execute the root cobra command
	cmd.Execute()
}
