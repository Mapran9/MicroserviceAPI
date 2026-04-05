package handlers

import (
	"database/sql"
	"log"

	"customer/internal/models"
	"customer/internal/repo"
	"github.com/gofiber/fiber/v2"
)

func GetCustomers(c *fiber.Ctx) error {
	customers, err := repo.GetAllCustomers()
	if err != nil {
		log.Printf("GetAllCustomers error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch customers",
		})
	}
	return c.JSON(customers)
}

func GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	customer, err := repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("GetCustomerByID error: %v", err)
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"message": "customer not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch customer",
		})
	}
	return c.JSON(customer)
}

func CreateCustomer(c *fiber.Ctx) error {
	var req models.Customer
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	newID, err := repo.CreateCustomerAutoID(c.Context(), &req)
	if err != nil {
		log.Printf("CreateCustomerAutoID error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to create customer",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "customer created successfully",
		"customer_id": newID,
	})
}
