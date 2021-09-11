package app

import (
	"log"
	"os"

	"github.com/arshabbir/consumer/controller"
)

func StartApplication() {
	//

	app := controller.NewEmpController()

	if app == nil {
		log.Println("Error Strting the application")
		os.Exit(1)
	}

	app.Start()
}
