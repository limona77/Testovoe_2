package controller

import (
	"testovoe_2/internal/service"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, services *service.Services) {
	auth := app.Group("/auth")
	newAuthRoutes(auth, services.IAuth)
	cR := &clientRoutes{clientService: services.IClient}
	app.Get("/me", cR.authMe)
	app.Post("/subscribe", cR.subscribe)
	app.Delete("/unsubscribe", cR.unSubscribe)
	app.Delete("/subscribe", cR.unSubscribe)
}
