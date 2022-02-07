package types

type ComposeFile struct {
	Version  string                 `yaml:"version"`
	Services map[string]interface{} `yaml:"services"`
}
