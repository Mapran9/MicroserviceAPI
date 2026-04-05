package routes

import (
	"github.com/gofiber/fiber/v2"

	"order/internal/handlers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	orders := api.Group("/Orders")
	orders.Post("/", handlers.CreateOrder)
	orders.Post("", handlers.CreateOrder)

	orders.Get("/:order_id", handlers.GetOrderByID)
	orders.Get("/customer/:customer_id", handlers.GetOrdersByCustomer)
}
