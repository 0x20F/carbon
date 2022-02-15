package types

type CarbonYaml struct {
	Name      string   `yaml:"-"`              // The name of the service
	Path      string   `yaml:"-"`              // The path where the carbon.yml was found
	Store     *Store   `yaml:"-"`              // The store this file belongs to
	Image     string   `yaml:"image"`          // The image of the service
	Container string   `yaml:"container_name"` // The container name of the service (will be overwritten by us)
	DependsOn []string `yaml:"depends_on"`     // The services that this service depends on

	// Everything within the file, unparsed
	FullContents ServiceFields
}

type CarbonConfig map[string]CarbonYaml

type ServiceDefinition map[string]ServiceFields

type ServiceFields map[string]interface{}
