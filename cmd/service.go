package cmd

import (
	"co2/builder"
	"co2/carbon"
	"co2/database"
	"co2/docker"
	"co2/helpers"
	"co2/runner"
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
		Run:   execService,
	}
)

func init() {
	serviceCmd.Flags().BoolVarP(&start, "start", "s", false, "whether to start the provided services")
	serviceCmd.Flags().BoolVarP(&stop, "stop", "p", false, "whether to stop the provided services")
}

func execService(cmd *cobra.Command, args []string) {
	services := strings.Join(args, " ")

	if start {
		_, ok := shouldRun(args)
		if !ok {
			return
		}

		fmt.Println("You want to START the following services:", services)
		handleStart(args)
	} else if stop {
		fmt.Println("You want to STOP the following services:", services)
		handleStop(args)
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

func handleStart(args []string) {
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

	// Build the command
	command := builder.DockerComposeCommand().
		File(compose.Path()).
		Service(strings.Join(args, " ")).
		Background().
		Up().
		Build()

	runner.Execute(command)

	go containerize(channel, compose)

	<-channel

}

// Handle the stopping of the provided running containers.
// This will make sure that each of the compose files will be found
// for the provided containers, and every single one of them will be
// stopped based on their individual compose files.
//
// While that's going on it will also make sure to remove the containers
// from the database since they are technically not running anymore.
func handleStop(args []string) {
	containers := database.Containers()
	grouped := map[string][]types.Container{}

	// Group the containers by their compose file
	for _, container := range containers {
		if helpers.Contains(args, container.Name) {
			grouped[container.ComposeFile] = append(grouped[container.ComposeFile], container)
		}
	}

	var wg sync.WaitGroup

	// Stop all the containers in each group and
	// delete them from the database
	for _, composeFile := range grouped {
		// Build the compose down command with all services
		// for each compose file
		command := builder.DockerComposeCommand().
			File(composeFile[0].ComposeFile).
			Down()

		for _, container := range composeFile {
			command.Service(container.Name)

			wg.Add(1)
			go func(container types.Container) {
				defer wg.Done()
				database.DeleteContainer(container)
			}(container)
		}

		// Run the command even if the database hasn't fully
		// updated yet. It's independent.
		runner.Execute(command.Build())
	}

	wg.Wait()
	fmt.Println("Stopped all required containers!")
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
