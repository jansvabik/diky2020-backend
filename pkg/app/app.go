package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// App is a structure of app data
type App struct {
	MongoClient *mongo.Client
	Cfg         config.Config
	Started     time.Time
}

// State contains the app data
var State App

// Create creates new app data structure
func Create() {
	State = App{}
}

// StatusHandler is a http handler for /
func StatusHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"name":    State.Cfg.App.Name,
		"env":     State.Cfg.App.Env,
		"started": State.Started,
		"links": map[string]string{
			"status":       "/",
			"homepageData": "/welcome",
			"thanks":       "/thanks",
		},
	})
}
