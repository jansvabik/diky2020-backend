package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/internal/thanks"
	"github.com/noltio/diky2020-backend/internal/welcome"
	"github.com/noltio/diky2020-backend/pkg/app"
)

// InitRoutes initialize all routes of this REST API service
func InitRoutes() {
	router := fiber.New()
	router.Get("/", app.StatusHandler)
	router.Get("/welcome", welcome.ReadHandler)
	router.Get("/thanks", thanks.ReadHandler)
	router.Get("/thanks/:page", thanks.ReadHandler)
	router.Post("/thanks", thanks.CreateHandler)
	router.Patch("/thanks/:id/donated", thanks.DonatedHandler)
	router.Listen(":" + strconv.Itoa(app.State.Cfg.Net.Port))
}
