package main

import (
	"co2/cmd"
	"math/rand"
	"time"
)

func main() {
	// Seed random for everything that needs it
	rand.Seed(time.Now().UnixNano())

	/*declarations := carbon.Configurations("../", 2)

	data, _ := yaml.Marshal(declarations["intake"].FullContents)
	fmt.Println(string(data))*/

	cmd.Execute()
	//runner.Execute("docker container ls -aq")
}
