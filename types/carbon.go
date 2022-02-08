package types

type CarbonYaml struct {
	Name      string
	Image     string   `yaml:"image"`
	Container string   `yaml:"container_name"`
	DependsOn []string `yaml:"depends_on"`

	// Everything within the file, unparsed
	FullContents ServiceFields
}

type CarbonConfig map[string]CarbonYaml

type ServiceDefinition map[string]ServiceFields

type ServiceFields map[string]interface{}
