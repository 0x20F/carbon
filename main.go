package main

import (
	"co2/database"
	"fmt"
)

func main() {
	/*declarations := carbon.Configurations("../", 2)

	data, _ := yaml.Marshal(declarations["intake"].FullContents)
	fmt.Println(string(data))*/

	containers := database.Containers()
	fmt.Println(containers)

}
