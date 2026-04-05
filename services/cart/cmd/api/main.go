package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"cart/config"
	"cart/internal/repo"
	"cart/internal/routes"
)

func main() {
	cfg := config.Load()

	repo.InitDB()

	app := fiber.New()
	routes.Setup(app)

	log.Printf("%s listening on :%s", cfg.ServiceName, cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
