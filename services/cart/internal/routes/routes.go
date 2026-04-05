package routes

import (
	"github.com/gofiber/fiber/v2"

	"cart/internal/handlers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	carts := api.Group("/Carts")
	carts.Post("/", handlers.AddCart)
	carts.Post("", handlers.AddCart)

	carts.Get("/customer/:customer_id", handlers.GetCartByCustomer)
	carts.Get("/:cart_id", handlers.GetCartByID)
	carts.Put("/:cart_id/status", handlers.UpdateCartStatus)
}
