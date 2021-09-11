package main

import (
	"log"

	"github.com/arshabbir/cassandraclient/app"
)

func main() {

	//os.Setenv("CLUSTERIP", "52.201.185.178")
	//os.Setenv("PORT", ":8080")
	log.Println("Starting the Application.......")
	app.StartApplication()

	return
}
