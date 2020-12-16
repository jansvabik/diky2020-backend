package router

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/internal/thanks"
	"github.com/noltio/diky2020-backend/internal/welcome"
	"github.com/noltio/diky2020-backend/pkg/app"
	"github.com/noltio/diky2020-backend/pkg/recaptcha"
)

// InitRoutes initialize all routes of this REST API service
func InitRoutes() {
	router := fiber.New()
	router.Get("/", app.StatusHandler)
	router.Get("/welcome", welcome.ReadHandler)
	router.Get("/thanks", thanks.ReadHandler)
	router.Get("/thanks/:page", thanks.ReadHandler)
	router.Post("/thanks", recaptcha.Middleware, thanks.CreateHandler)
	router.Post("/thanks/:id/donation", thanks.DonatedHandler)
	router.Post("/thanks/:id/likes", thanks.LikeHandler)
	router.Listen(":" + strconv.Itoa(app.State.Cfg.Net.Port))
}
