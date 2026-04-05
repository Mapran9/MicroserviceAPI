package main

import (
	"log"

	"customer/config"
	"customer/internal/repo"
	"customer/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	repo.InitDB()

	app := fiber.New()
	routes.Setup(app)

	log.Printf("%s listening on :%s", cfg.ServiceName, cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
