package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/internal/thanks"
	"github.com/noltio/diky2020-backend/pkg/app"
)

// InitRoutes initialize all routes of this REST API service
func InitRoutes() {
	router := fiber.New()
	router.Get("/", app.StatusHandler)
	router.Get("/thanks", thanks.ReadHandler)
	router.Post("/thanks", thanks.CreateHandler)
	router.Listen(":80")
}
