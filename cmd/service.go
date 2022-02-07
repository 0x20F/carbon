package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	start bool
	stop  bool

	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Does things with the services you provide",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			services := strings.Join(args, " ")

			// 0. Create a new docker file name
			// 1. For each of the provided services
			// 2. Get their configuration
			// 3. Inject a random container name so nothing clashes
			// 4. Group them all together in a map
			// 5. Save that Map into a docker compose file
			// 6. for now, output the new file.

			if start && !stop {
				fmt.Println("You want to START the following services:", services)
			}

			if stop && !start {
				fmt.Println("You want to STOP the following services:", services)
			}
		},
	}
)

func init() {
	serviceCmd.Flags().BoolVarP(&start, "start", "s", false, "whether to start the provided services")
	serviceCmd.Flags().BoolVarP(&stop, "stop", "p", false, "whether to stop the provided services")
}
