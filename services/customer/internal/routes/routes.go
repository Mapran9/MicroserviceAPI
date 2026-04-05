package routes

import (
	"customer/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	customer := api.Group("/customers")

	customer.Get("", handlers.GetCustomers)
	customer.Get("/", handlers.GetCustomers)

	customer.Get("/:id", handlers.GetCustomer)

	customer.Post("", handlers.CreateCustomer)
	customer.Post("/", handlers.CreateCustomer)
}
