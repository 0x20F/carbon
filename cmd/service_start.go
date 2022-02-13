package cmd

import (
	"co2/builder"
	"co2/carbon"
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	force bool

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the provided services",
		Args:  cobra.MinimumNArgs(1),
		Run:   start,
	}
)

func init() {
	help := "Force the start of the service. This will delete the old ones before starting."
	startCmd.Flags().BoolVarP(&force, "force", "f", false, help)
}

// Checks whether or not the containers are already
// in the database. If they are, we don't want to
// start them again.
//
// If the force flag is provided, this will always return
// true./
func shouldRun(provided []string) bool {
	if force {
		return true
	}

	// Get all containers from the database
	containers := database.Containers()

	// If an of the provided containers is in the database, quit
	for _, container := range containers {
		if helpers.Contains(provided, container.Name) {
			printer.Error("ERROR", "service already running:", container.Name)
			printer.Extra(
				printer.Red,
				"Try running `co2 stop` first",
				"Use the `--force` flag to force the start",
			)
			return false
		}
	}

	return true
}

// Looks through all the registered stores and returns all
// the carbon services that are defined within those stores.
func services() types.CarbonConfig {
	printer.Extra(printer.Green, "Looking through the store")

	stores := database.Stores()
	configs := types.CarbonConfig{}

	// For each store
	for _, store := range stores {
		// Find all carbon files in the store
		files := carbon.Configurations(store.Path, 2)

		// For each carbon file
		for k, v := range files {
			configs[k] = v
		}
	}

	return configs
}

// Looks through all the available services and returns only
// the ones that the user has specified in the command.
//
// If no file is found for a specific service, it will output
// some information to stdout and continue to the next one.
func extract(args []string) types.CarbonConfig {
	choices := types.CarbonConfig{}
	configs := services()

	for _, service := range args {
		if _, ok := configs[service]; !ok {
			printer.Extra(printer.Red, "No carbon file found for: "+service)
			continue
		}

		found := configs[service]
		if _, ok := found.FullContents["depends_on"]; ok {
			omitted := false

			for _, dep := range found.FullContents["depends_on"].([]interface{}) {
				if helpers.Contains(args, dep.(string)) {
					continue
				}

				message := fmt.Sprintf("'%s' depends on '%s' but '%s' is not provided", service, dep.(string), dep.(string))
				printer.Extra(printer.Cyan, message)

				omitted = true
			}

			if omitted {
				continue
			}
		}

		container := service + "-" + helpers.RandomAlphaString(10)
		found.FullContents["container_name"] = container
		found.Container = container
		choices[service] = found
	}

	return choices
}

// Generates a new compose file structure based on the provided
// services, if they exist. This will make sure to inject all of
// the required values into all the containers within the compose
// file.
func compose(args []string) (types.ComposeFile, error) {
	choices := extract(args)
	if len(choices) == 0 {
		return types.ComposeFile{}, errors.New("no services found")
	}

	printer.Extra(printer.Green, "Generating compose file")
	compose := types.NewComposeFile()

	// Add all the services to the compose file
	for _, service := range choices {
		compose.Services[service.Name] = service.FullContents
	}

	printer.Extra(printer.Green, "Saving compose file to `"+compose.Path()+"`")
	compose.Save()

	// Save all containers to the database
	channel := make(chan bool)
	go containerize(channel, compose)
	<-channel

	return compose, nil
}

// Creates container types for each of the provided services
// within the compose file, making sure that all the containers
// know which file they belong to.
//
// Also saves all the containers to the database so that all required
// information can be retrieved later if ever needed.
func containerize(channel chan bool, compose types.ComposeFile) {
	// Create container types for all services in the compose file
	containers := []types.Container{}

	for name, service := range compose.Services {
		container := types.Container{
			Name:        name,
			Uid:         service["container_name"].(string),
			ComposeFile: compose.Path(),
		}

		containers = append(containers, container)
	}

	var wg sync.WaitGroup

	// Save all the containers
	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()

			// Save the container
			database.AddContainer(container)
		}(container)
	}

	wg.Wait()
	channel <- true
}

// Starts the service start command.
// This is where all of the logging of the command happens
// as well.
//
// If the force flag is provided, this will make sure to call
// the stop command beforehand so that all the services we're
// trying to start will start fresh.
func start(cmd *cobra.Command, args []string) {
	if ok := shouldRun(args); !ok {
		return
	}

	printer.Info(
		printer.Green,
		"START",
		"Starting provided services:",
		strings.Join(args, ", "),
	)

	// Before anything else, make sure we find
	// all the services we require for the start.
	file, err := compose(args)
	if err != nil {
		printer.Extra(printer.Grey, "Aborting")
		return
	}

	// If we're forcing, we want to cleanup first
	if force {
		printer.Extra(
			printer.Yellow,
			"`--force` flag is set, stopping all provided services first",
		)

		execStop(cmd, args)
	}

	// Execute the compose command
	printer.Extra(printer.Green, "Executing `docker compose` command on the new file\n")
	command := builder.DockerComposeCommand().
		File(file.Path()).
		Service(strings.Join(args, " ")).
		Background().
		Up().
		Build()

	runner.Execute(command)
}
