package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noltio/tiskarnik-ms-auth/pkg/app"
)

// InitRoutes initialize all routes of this REST API service
func InitRoutes() {
	router := fiber.New()
	router.Get("/", app.StatusHandler)
	router.Listen(":80")
}
