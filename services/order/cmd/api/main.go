package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"order/config"
	"order/internal/repo"
	"order/internal/routes"
)

func main() {
	cfg := config.Load()

	repo.InitDB()

	app := fiber.New()
	app.Use(recover.New())

	routes.Setup(app)

	log.Printf("%s listening on :%s instance=%s", cfg.ServiceName, cfg.Port, cfg.InstanceID)
	log.Fatal(app.Listen(":" + cfg.Port))
}
