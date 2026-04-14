package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"payment/internal/handlers"
)

func Setup(app *fiber.App) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"service":  getEnv("SERVICE_NAME", "payment-service"),
			"instance": getInstanceID(),
		})
	})

	app.Get("/whoami", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":  getEnv("SERVICE_NAME", "payment-service"),
			"instance": getInstanceID(),
			"hostname": getEnv("HOSTNAME", "unknown"),
		})
	})

	api := app.Group("/api")

	payments := api.Group("/payments")
	payments.Post("/internal", handlers.CreatePayment)
}

func getInstanceID() string {
	if v := os.Getenv("INSTANCE_ID"); v != "" {
		return v
	}
	if v := os.Getenv("HOSTNAME"); v != "" {
		return v
	}
	return "unknown"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
