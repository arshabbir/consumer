package app

import (
	"log"
	"os"

	"github.com/arshabbir/cassandraclient/controller"
)

func StartApplication() {
	//

	app := controller.NewStudentController()

	if app == nil {
		log.Println("Error Strting the application")
		os.Exit(1)
	}

	app.Start()
}
