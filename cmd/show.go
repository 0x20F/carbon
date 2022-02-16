package cmd

import (
	"co2/database"
	"co2/docker"
	"co2/printer"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	running   bool
	stores    bool
	available bool

	showCmd = &cobra.Command{
		Use:   "show",
		Short: "Shows different kinds of information",
		Run:   execShow,
	}

	// Some data in the displayed tables might not be
	// that important so we want it to have a faded color.
	// Still visible, but less important.
	fadedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))
)

// Adds all the required flags
func init() {
	showCmd.Flags().BoolVarP(&running, "running", "r", false, "show all currently running containers")
	showCmd.Flags().BoolVarP(&stores, "stores", "s", false, "show all registered stores")
	showCmd.Flags().BoolVarP(&available, "carbon", "c", false, "show all available carbon services")
}

// Checks what flags are provided and displays
// the specific table representing each flag.
func execShow(cmd *cobra.Command, args []string) {
	if running {
		showRunning()
	}

	if stores {
		showStores()
	}

	if available {
		showAvailable()
	}
}

// Shows all the running containers.
//
// Since this runs using the docker api, it's completely
// independent of any carbon features, or even the
// database so it's pretty fast.
//
// The resulting table should also contain the unique
// id for the container generated from the name and the image.
func showRunning() {
	containers := docker.RunningContainers()

	if len(containers) == 0 {
		printer.Info(printer.Cyan, "RUN", "No running containers", "")
		return
	}

	table := printer.NewTable(7)
	printer.Info(printer.Cyan, "RUN", "total running containers:", fmt.Sprint(len(containers)))

	table.Header(
		"KEY",
		"NAME",
		"ID",
		"IMAGE",
		"PORTS",
		"CREATED",
		"STATUS",
	)

	for key, container := range containers {
		ports := []string{}

		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", port.PublicPort, port.Type))
		}

		table.Row(
			key,
			container.Names[0],
			container.ID[:7],
			container.Image,
			fadedStyle.Render(strings.Join(ports, ", ")), // Turn the array of ports into a string
			fadedStyle.Render(fmt.Sprint(container.Created)),
			fadedStyle.Render(container.Status),
		)
	}

	table.Display()
}

// Shows all the registered carbon stores.
//
// It's not that complicated, thankfully.
// Just a simple query to the database to fetch all
// the currently registered stores and then display it
// in a nicely formatted table.
//
// If the stores don't have an environment file set,
// this will replace the value with 'undefined' in the
// resulting table.
func showStores() {
	stores := database.Stores()

	if len(stores) == 0 {
		printer.Info(printer.Grey, "STORE", "No registered stores", "")
		return
	}

	table := printer.NewTable(4)
	printer.Info(printer.Grey, "STORE", "total registered stores:", fmt.Sprint(len(stores)))

	table.Header(
		"KEY",
		"PATH",
		"DATE",
		"ENV",
	)

	for _, store := range stores {
		env := "undefined"

		if store.Env != "" {
			env = store.Env
		}

		table.Row(
			store.Uid,
			store.Path,
			fadedStyle.Render(fmt.Sprint(store.CreatedAt)),
			env,
		)
	}

	table.Display()
}

// Shows all the available carbon services.
//
// This will display a table of all the defined carbon
// services (carbon.yml) by looking through all of the stores,
// and then looking through all the directories within each store.
//
// All the long paths are shortened so they don't occupy too
// much screen space.
func showAvailable() {
	services := services()

	if len(services) == 0 {
		printer.Info(printer.Grey, "CARBON", "No available carbon services", "")
		return
	}

	keys := make([]string, 0, len(services))

	for key := range services {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	table := printer.NewTable(3)
	printer.Info(printer.Grey, "CARBON", "total available carbon services:", fmt.Sprint(len(services)))

	table.Header(
		"NAME",
		"IMAGE",
		"PATH",
	)

	for _, name := range keys {
		service := services[name]

		table.Row(
			name,
			service.Image,
			fadedStyle.Render(fmt.Sprintf("...%s", service.Path[len(service.Path)-30:])),
		)
	}

	table.Display()
}
