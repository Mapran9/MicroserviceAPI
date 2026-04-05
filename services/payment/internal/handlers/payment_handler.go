package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"payment/internal/repo"
)

type CreatePaymentRequest struct {
	CustomerID    string  `json:"customer_id"`
	OrderID       string  `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

type CreatePaymentResponse struct {
	Message   string  `json:"message"`
	PaymentID string  `json:"payment_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
}

func CreatePayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	if req.CustomerID == "" || req.OrderID == "" || req.PaymentMethod == "" {
		return c.Status(400).JSON(fiber.Map{"message": "customer_id, order_id and payment_method are required"})
	}

	if req.Status == "" {
		req.Status = "pending"
	}

	paymentID, err := repo.CreatePaymentAutoID(
		c.Context(),
		req.CustomerID,
		req.OrderID,
		req.PaymentMethod,
		req.Amount,
		req.Status,
	)
	if err != nil {
		log.Printf("CreatePaymentAutoID error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to create payment"})
	}

	return c.Status(201).JSON(CreatePaymentResponse{
		Message:   "Payment created successfully",
		PaymentID: paymentID,
		Amount:    req.Amount,
		Status:    req.Status,
	})
}
