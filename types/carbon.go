package types

// Single service definition for a carbon.yml file.
// This is what we care aboout from the things that
// a user writes in a carbon configuration file.
type CarbonService struct {
	Name      string   `yaml:"-"`              // The name of the service
	Path      string   `yaml:"-"`              // The path where the carbon.yml was found
	Store     *Store   `yaml:"-"`              // The store this file belongs to
	Image     string   `yaml:"image"`          // The image of the service
	Container string   `yaml:"container_name"` // The container name of the service (will be overwritten by us)
	DependsOn []string `yaml:"depends_on"`     // The services that this service depends on

	// Everything within the file, unparsed
	FullContents ServiceFields
}

// Alias type for a map of carbon services
type CarbonConfig map[string]CarbonService

// Alias type for a map of map of unknown
type ServiceDefinition map[string]ServiceFields

// Alias type for a map of unknown
type ServiceFields map[string]interface{}
