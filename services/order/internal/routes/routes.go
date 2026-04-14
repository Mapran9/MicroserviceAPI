package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"order/internal/handlers"
)

func Setup(app *fiber.App) {
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"service":  getEnv("SERVICE_NAME", "order-service"),
			"instance": getInstanceID(),
		})
	})

	app.Get("/whoami", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service":  getEnv("SERVICE_NAME", "order-service"),
			"instance": getInstanceID(),
			"hostname": getEnv("HOSTNAME", "unknown"),
		})
	})

	api := app.Group("/api")

	orders := api.Group("/Orders")
	orders.Post("/", handlers.CreateOrder)
	orders.Post("", handlers.CreateOrder)

	orders.Get("/:order_id", handlers.GetOrderByID)
	orders.Get("/customer/:customer_id", handlers.GetOrdersByCustomer)
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
