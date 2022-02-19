package cmd

import (
	"co2/builder"
	"co2/database"
	"co2/helpers"
	"co2/printer"
	"co2/runner"
	"co2/types"
	"errors"
	"fmt"
	"strings"

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

// Adds all the required flags
func init() {
	help := "Force the start of the service. This will delete the old ones before starting."
	startCmd.Flags().BoolVarP(&force, "force", "f", false, help)
}

// Starts the service start command.
// This is where all of the logging of the command happens
// as well.
//
// If the force flag is provided, this will make sure to call
// the stop command beforehand so that all the services we're
// trying to start will start fresh.
//
// We also want to make sure that we tell the docker compose command
// to run with any of the available environment files that might be
// provided by the current store we are looking at.
func start(cmd *cobra.Command, args []string) {
	if ok := shouldRun(args, force); !ok {
		return
	}

	printer.Info(
		printer.Green,
		"START",
		"Starting provided services:",
		strings.Join(args, ", "),
	)

	// If we're forcing, we want to cleanup first
	if force {
		printer.Extra(
			printer.Yellow,
			"`--force` flag is set, stopping all provided services first",
		)

		execStop(cmd, args)
	}

	extracted := extract(args)
	envs, composeFile, err := compose(extracted)
	if err != nil {
		printer.Extra(printer.Grey, "Aborting")
		return
	}
	containerize(composeFile)
	run(composeFile, envs, args)
}

// Generates and runs the docker compose command based on the
// resolved compose file, environment files, and the services
// that the user has provided.
func run(file types.ComposeFile, envs []string, services []string) {
	command := builder.DockerComposeCommand().
		File(file.Path()).
		Service(strings.Join(services, " ")).
		Background().
		Up()

	// Add all the environment files as well
	for _, env := range envs {
		command.EnvFile(env)
	}

	printer.Extra(printer.Green, "Executing `docker compose` command on the new file\n")
	runner.Execute(types.Command{
		Text: command.Build(),
	})
}

// Checks whether or not the containers are already
// in the database. If they are, we don't want to
// start them again.
//
// If the force flag is provided, this will always return
// true.
func shouldRun(choices []string, force bool) bool {
	if force {
		return true
	}

	// Get all containers from the database
	containers := database.Containers()

	// If an of the provided containers is in the database, quit
	for _, container := range containers {
		if helpers.Contains(choices, container.ServiceName) {
			printer.Error("ERROR", "service already running:", container.ServiceName)
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

// Looks through all the available services and returns only
// the ones that the user has specified in the command.
//
// If no file is found for a specific service, it will output
// some information to stdout and continue to the next one.
//
// If any of the user provided services are dependent on other
// services but the other services aren't already provided, this will
// inform the user about their error and ignore the service and move
// onto the next one.
func extract(args []string) types.CarbonConfig {
	printer.Extra(printer.Green, "Looking through the store")

	choices := types.CarbonConfig{}
	configs := fs.Services()

	for _, service := range args {
		if _, ok := configs[service]; !ok {
			printer.Extra(printer.Red, "No carbon file found for: "+service)
			continue
		}

		found := configs[service]
		omitted := false

		for _, dep := range found.DependsOn {
			if helpers.Contains(args, dep) {
				continue
			}

			message := fmt.Sprintf("'%s' depends on '%s' but '%s' is not provided, ignoring.", service, dep, dep)
			printer.Extra(printer.Cyan, message)

			omitted = true
		}

		if omitted {
			continue
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
func compose(choices types.CarbonConfig) ([]string, types.ComposeFile, error) {
	envs := []string{}
	if len(choices) == 0 {
		return envs, types.ComposeFile{}, errors.New("no services found")
	}

	printer.Extra(printer.Green, "Generating compose file")
	compose := types.NewComposeFile()

	// Add all the services to the compose file
	for _, service := range choices {
		compose.Services[service.Name] = service.FullContents
	}

	printer.Extra(printer.Green, "Saving compose file to `"+compose.Path()+"`")
	compose.Save()

	// Find all the env files that should be given to the compose file
	for _, service := range choices {
		if service.Store.Env == "" {
			continue
		}

		if !helpers.Contains(envs, service.Store.Env) {
			envs = append(envs, service.Store.Env)
		}
	}

	return envs, compose, nil
}

// Creates container types for each of the provided services
// within the compose file, making sure that all the containers
// know which file they belong to.
//
// Also saves all the containers to the database so that all required
// information can be retrieved later if ever needed.
func containerize(compose types.ComposeFile) {
	containers := []types.Container{}

	for name, service := range compose.Services {
		container := types.Container{
			ServiceName: name,
			Name:        service["container_name"].(string),
			Image:       service["image"].(string),
			Status:      "Created",
			ComposeFile: compose.Path(),
		}
		container.Hash()

		containers = append(containers, container)
	}

	// Try saving it all async so it goes faster,
	// we can do other things in the meantime if we ever need to.
	for _, container := range containers {
		database.AddContainer(container)
	}
}
