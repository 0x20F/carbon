package cmd

import (
	"co2/carbon"
	"co2/database"
	"co2/docker"
	"co2/helpers"
	"co2/types"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	start bool
	stop  bool

	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Does things with the services you provide",
		Args:  cobra.MinimumNArgs(1),
		Run:   exec,
	}
)

func init() {
	serviceCmd.Flags().BoolVarP(&start, "start", "s", false, "whether to start the provided services")
	serviceCmd.Flags().BoolVarP(&stop, "stop", "p", false, "whether to stop the provided services")
}

func exec(cmd *cobra.Command, args []string) {
	_, ok := shouldRun(args)
	if !ok {
		return
	}

	services := strings.Join(args, " ")

	configs := carbon.Configurations("../", 2) // FIXME: Don't hardcode this
	choices := types.CarbonConfig{}

	// 2. Get their configurations
	for _, service := range args {
		if _, ok := configs[service]; !ok {
			continue
		}

		found := configs[service]

		// Inject a random container name
		container := service + "-" + helpers.RandomAlphaString(10)
		found.FullContents["container_name"] = container
		found.Container = container

		choices[service] = found
	}

	compose := docker.NewComposeFile()

	// Add all the services to the compose file
	for _, service := range choices {
		compose.Services[service.Name] = service.FullContents
	}

	// Save the compose file
	channel := make(chan bool)
	go compose.Save(channel)

	<-channel

	go containerize(channel, compose)

	<-channel

	if start && !stop {
		fmt.Println("You want to START the following services:", services)
	}

	if stop && !start {
		fmt.Println("You want to STOP the following services:", services)
	}
}

func shouldRun(provided []string) ([]types.Container, bool) {
	// Get all containers from the database
	containers := database.Containers()

	// If an of the provided containers is in the database, quit
	for _, container := range containers {
		if helpers.Contains(provided, container.Name) {
			fmt.Printf("%s is already in the database\n", container.Name)
			return nil, false
		}
	}

	return containers, true
}

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

	// Create a workgroup for all the containers
	var wg sync.WaitGroup

	// Save all the containers
	for _, container := range containers {
		wg.Add(1)
		go func(container types.Container) {
			defer wg.Done()

			// Save the container
			database.InsertContainer(container)
		}(container)
	}

	wg.Wait()
	channel <- true
}
