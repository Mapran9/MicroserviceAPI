package routes

import (
	"github.com/gofiber/fiber/v2"

	"payment/internal/handlers"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	payments := api.Group("/payments")
	payments.Post("/internal", handlers.CreatePayment)
}
