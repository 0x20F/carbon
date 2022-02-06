package main

import (
	"co2/database"
	"fmt"
)

func main() {
	/*declarations := carbon.Configurations("../", 2)

	data, _ := yaml.Marshal(declarations["intake"].FullContents)
	fmt.Println(string(data))*/

	// Next Steps:
	// 1. User provides list of carbon defined services to the program
	// 2. Each of those services will get an injected container_name
	// 3. All relevant services get merged into their own docker compose file
	// 4. That docker compose file gets saved under a random name
	// 5. Docker compose file gets started up
	// 6. All service data gets pushed to the DB if compose file has started without issues

	containers := database.Containers()
	fmt.Println(containers)

}
