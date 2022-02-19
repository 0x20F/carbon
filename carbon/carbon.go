package carbon

import (
	"bytes"
	"co2/types"
	"io/ioutil"
	"reflect"

	"gopkg.in/yaml.v2"
)

// Recursively look through a directory until the given
// depth is reached.
//
// If, along the way, we find a carbon.yml file, we
// will parse it and return a map with all the found
// services within the file.
//
// If we find a directory and the max depth hasn't been reached
// yet, we go deeper and look for carbon.yml files.
func findCarbonFiles(root string, depth int) []string {
	// Break if the path isn't a directory
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
			// Make sure we get the full path of the carbon.yml file
			parsed = append(parsed, root+"/"+file.Name())
		}
	}

	return parsed
}

// This will take all the existing carbon configurations and
// make sure to map them into usable data structures.
//
// Each service definition will be mapped into a CarbonYaml structure
// which in turn will be mapped along with other yaml structures into a
// CarbonConfig structure.
//
// This also injects all the required data within all the CarbonYaml
// structures. Such as the Path where it was found, and the full contents
// of service it belongs to so it's easy to remap all of it into a new
// yaml file later on.
func Configurations(path string, depth int) types.CarbonConfig {
	files := findCarbonFiles(path, depth)

	var config types.CarbonConfig = make(types.CarbonConfig, len(files))

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		documents := documents(content, file)

		for k, v := range documents {
			config[k] = v
		}
	}

	return config
}

// Separates the given file contents at the yaml document
// separator and parses all the found documents according to
// the carbon requirements.
//
// Each service will be parsed into a CarbonConfig, and then into
// an arbitrary map that contains all the fields of the defined configuration.
// This arbitrary map will then be mapped into a full CarbonConfig so that
// each service has access to all of the contents within their service definition
// if they ever need it.
func documents(contents []byte, file string) types.CarbonConfig {
	documents := bytes.Split(contents, []byte("---"))

	var final types.CarbonConfig = make(types.CarbonConfig)

	for _, doc := range documents {
		full := types.CarbonConfig{}
		fake := types.ServiceDefinition{}

		// Unmarshal once into a structure with limited fields
		err := yaml.Unmarshal(doc, &full)
		if err != nil {
			panic(err)
		}

		// Unmarshal again into an arbitrary map with all the available fields
		err = yaml.Unmarshal(doc, &fake)
		if err != nil {
			panic(err)
		}

		// Map the values from the fake map into the real map
		k, v := move(fake, full, file)
		final[k] = v
	}

	return final
}

// Maps the required fields from the arbitrary map with no specific structure
// into the real CarbonConfig map with the correct structure.
//
// This makes sure that the final service representation knows the path
// it came from, the name of the service it represents, and has access
// to all of the contents within their service definition if they ever
// need it.
//
// This has to exist since we're not building a 1:1 mapping between a
// docker-compose service definition but we still want all the data.
func move(this types.ServiceDefinition, into types.CarbonConfig, file string) (string, types.CarbonService) {
	// Get all the key from the map even though we know there's only 1
	key := reflect.ValueOf(into).MapKeys()[0].String()
	current := into[key]

	current.Path = file
	current.Name = key
	current.FullContents = make(types.ServiceFields)
	current.FullContents = this[key]

	return key, current
}
