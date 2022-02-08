package cmd

import (
	"co2/carbon"
	"co2/docker"
	"co2/helpers"
	"co2/types"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
		container := helpers.RandomAlphaString(10)
		found.FullContents["container_name"] = container
		found.Container = container

		choices[service] = found
	}

	compose := docker.NewComposeFile()

	// Add all the services to the compose file
	for _, service := range choices {
		compose.Services[service.Name] = service.FullContents
	}

	// Turn it into yaml
	yml, err := yaml.Marshal(compose)
	if err != nil {
		panic(err)
	}

	// Print the yaml
	fmt.Println(string(yml))

	// 4. Group them all together in a map
	// 5. Save that Map into a docker compose file
	// 6. for now, output the new file.

	if start && !stop {
		fmt.Println("You want to START the following services:", services)
	}

	if stop && !start {
		fmt.Println("You want to STOP the following services:", services)
	}
}
