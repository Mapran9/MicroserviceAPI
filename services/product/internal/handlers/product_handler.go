package handlers

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"

	"product/internal/models"
	"product/internal/repo"
)

func GetProducts(c *fiber.Ctx) error {
	products, err := repo.GetAllProducts()
	if err != nil {
		log.Printf("GetAllProducts error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch products"})
	}
	return c.JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := repo.GetProductByID(id)
	if err != nil {
		log.Printf("GetProductByID error: %v", err)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": "product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch product"})
	}
	return c.JSON(product)
}

// POST /api/Products -> auto-generate product_id = P00001...
func CreateProduct(c *fiber.Ctx) error {
	var req models.Product
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	newID, err := repo.CreateProductAutoID(c.Context(), &req)
	if err != nil {
		log.Printf("CreateProductAutoID error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to create product",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":    "product created successfully",
		"product_id": newID,
	})
}
