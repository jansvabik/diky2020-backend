package main

import (
	"log"
	"os"
	"time"

	"github.com/noltio/tiskarnik-ms-auth/pkg/app"
	"github.com/noltio/tiskarnik-ms-auth/pkg/config"
	"github.com/noltio/tiskarnik-ms-auth/pkg/db"
	"github.com/noltio/tiskarnik-ms-auth/pkg/router"
)

func main() {
	// create an app and store its configuration
	app.Create()
	err := config.Load(&app.State.Cfg)
	if err != nil {
		log.Fatalln("Failed to load the configuration")
		os.Exit(1)
	}

	// try to connect to the database
	client, err := db.MongoConnect()
	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
		os.Exit(1)
	}
	app.State.MongoClient = client

	// save start time
	app.State.Started = time.Now()

	// init routers
	router.InitRoutes()
}
