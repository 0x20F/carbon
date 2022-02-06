package main

import (
	"co2/carbon"
	"fmt"

	"gopkg.in/yaml.v2"
)

func main() {
	declarations := carbon.Configurations("../", 2)

	data, _ := yaml.Marshal(declarations["intake"].FullContents)

	fmt.Println(string(data))
}
