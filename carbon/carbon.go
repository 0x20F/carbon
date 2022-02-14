package carbon

import (
	"bytes"
	"co2/types"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
func Configurations(path string, depth int) types.CarbonConfig {
	files := findCarbonFiles(path, depth)

	// We want one map with all the available
	// configurations for a given path.
	var config types.CarbonConfig = make(types.CarbonConfig)
	var values types.ServiceDefinition = make(types.ServiceDefinition)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		// Split at yaml document separator
		documents := bytes.Split(content, []byte("---"))

		for _, doc := range documents {
			temp := types.CarbonConfig{}

			err := yaml.Unmarshal(doc, &temp)
			if err != nil {
				panic(err)
			}

			// Inject the path into the config
			for k, v := range temp {
				v.Path = file
				config[k] = v
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

		c.Name = key
		c.FullContents = make(types.ServiceFields)
		c.FullContents = v

		config[key] = c

		count++
	}

	return config
}
