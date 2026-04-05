package routes

import (
	"github.com/gofiber/fiber/v2"

	"product/internal/handlers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	product := api.Group("/Products")
	product.Get("", handlers.GetProducts)
	product.Get("/", handlers.GetProducts)
	product.Get("/:id", handlers.GetProduct)

	product.Post("", handlers.CreateProduct)
	product.Post("/", handlers.CreateProduct)
}
