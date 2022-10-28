package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/radekkrejcirik01/Casblanca-backend/services/login-registration/pkg/rest/controller"
)

// Create new REST API serveer
func Create() *fiber.App {
	app := fiber.New()

	app.Get("/", controller.Index)

	app.Post("/register", controller.UserRegister)
	app.Post("/login", controller.UserLogin)
	app.Post("/tags", controller.AddTag)
	app.Post("/photos", controller.AddTag)
	app.Get("/get", controller.UserGet)
	app.Put("/update", controller.UserPut)
	app.Delete("/delete", controller.UserDelete)

	return app
}
