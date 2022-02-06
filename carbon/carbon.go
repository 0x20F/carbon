package carbon

import (
	"bytes"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type CarbonYaml struct {
	Image     string   `yaml:"image"`
	Container string   `yaml:"container_name"`
	DependsOn []string `yaml:"depends_on"`

	// Everything within the file, unparsed
	FullContents map[string]interface{}
}

type CarbonConfig map[string]CarbonYaml

// Recursively dig through a directory until
// we exhaust all the available directories in the given
// path while making sure to keep track of all the
// carbon.yml files.
func findCarbonFiles(root string, depth int) []string {
	paths, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	parsed := []string{}

	for _, file := range paths {
		if depth == 0 {
			break
		}

		if file.IsDir() {
			fresh := findCarbonFiles(root+"/"+file.Name(), depth-1)
			parsed = append(parsed, fresh...)
			continue
		}

		if file.Name() == "carbon.yml" {
			parsed = append(parsed, root+"/"+file.Name())
		}
	}

	return parsed
}

// Look through a specific directory and find all the
// carbon.yml files. Parse them and return a nice map with
// all the found services.
func Configurations(path string, depth int) CarbonConfig {
	files := findCarbonFiles(path, depth)

	// We want one map with all the available
	// configurations for a given path.
	var config CarbonConfig = make(CarbonConfig)
	var values map[string]interface{} = make(map[string]interface{})

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		// Split at yaml document separator
		documents := bytes.Split(content, []byte("---"))

		for _, doc := range documents {
			err := yaml.Unmarshal(doc, &config)
			if err != nil {
				panic(err)
			}

			err = yaml.Unmarshal(doc, &values)
			if err != nil {
				panic(err)
			}
		}
	}

	// Map all unmarshaled file contents to
	// each of the generated Carbon configuration files.
	count := 0
	for key := range config {
		c := config[key]
		v := values[key]

		c.FullContents = make(map[string]interface{})
		c.FullContents[key] = v

		config[key] = c

		count++
	}

	return config
}
