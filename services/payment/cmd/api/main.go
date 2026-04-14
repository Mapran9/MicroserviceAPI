package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"payment/config"
	"payment/internal/repo"
	"payment/internal/routes"
)

func main() {
	cfg := config.Load()

	repo.InitDB()

	app := fiber.New()
	routes.Setup(app)

	log.Printf("%s listening on :%s instance=%s", cfg.ServiceName, cfg.Port, cfg.InstanceID)
	log.Fatal(app.Listen(":" + cfg.Port))
}
