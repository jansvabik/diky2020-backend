package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/tiskarnik-ms-auth/pkg/config"
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
	return c.JSON(struct {
		Name    string    `json:"name"`
		Env     string    `json:"env"`
		Started time.Time `json:"started"`
	}{
		Name:    State.Cfg.App.Name,
		Env:     State.Cfg.App.Env,
		Started: State.Started,
	})
}
