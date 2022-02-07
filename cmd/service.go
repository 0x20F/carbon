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
