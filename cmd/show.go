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

	fadedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777"))
)

func init() {
	showCmd.Flags().BoolVarP(&running, "running", "r", false, "show all currently running containers")
	showCmd.Flags().BoolVarP(&stores, "stores", "s", false, "show all registered stores")
	showCmd.Flags().BoolVarP(&available, "carbon", "c", false, "show all available carbon services")
}

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
		// Turn the array of ports into a string
		ports := []string{}

		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%d/%s", port.PublicPort, port.Type))
		}

		table.Row(
			key,
			container.Names[0],
			container.ID[:7],
			container.Image,
			fadedStyle.Render(strings.Join(ports, ", ")),
			fadedStyle.Render(fmt.Sprint(container.Created)),
			fadedStyle.Render(container.Status),
		)
	}

	table.Display()
}

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
