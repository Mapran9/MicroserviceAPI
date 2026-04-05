package handlers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
	_ "strings"

	"cart/internal/clients"
	"cart/internal/repo"
)

// ===== Request/Response DTO (match Mono API doc) =====

type AddCartRequest struct {
	CustomerID string `json:"customer_id"`
	Items      []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

type AddCartResponse struct {
	CartID  string `json:"cart_id"`
	Message string `json:"message"`
}

type CartItemResponse struct {
	CartItemID string  `json:"cart_item_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type CartResponse struct {
	CartID     string             `json:"cart_id"`
	CustomerID string             `json:"customer_id"`
	Status     string             `json:"status"`
	Items      []CartItemResponse `json:"items"`
}

// POST /api/Carts
// - สร้าง cart ใหม่ทุกครั้ง (ตามตัวอย่างในเอกสารที่สร้าง CART000010 แล้วสร้าง CART000011 ได้อีก) :contentReference[oaicite:5]{index=5}
// - status = "pending" :contentReference[oaicite:6]{index=6}
// - ดึง price จาก product-service แล้วค่อย insert cart_items
func AddCart(c *fiber.Ctx) error {
	var req AddCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	if req.CustomerID == "" || len(req.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "customer_id and items are required"})
	}

	// 1) Create cart header (new cart every request) with status "pending"
	cartID, err := repo.CreateCartAutoID(c.Context(), req.CustomerID, "pending")
	if err != nil {
		log.Printf("CreateCartAutoID error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to create cart"})
	}

	// 2) Add items (fetch price from product-service)
	for _, it := range req.Items {
		if it.ProductID == "" || it.Quantity <= 0 {
			return c.Status(400).JSON(fiber.Map{"message": "each item must have product_id and quantity > 0"})
		}

		price, err := clients.GetProductPrice(it.ProductID)
		if err != nil {
			log.Printf("GetProductPrice error (product_id=%s): %v", it.ProductID, err)
			return c.Status(400).JSON(fiber.Map{"message": "invalid product_id: " + it.ProductID})
		}

		if _, err := repo.AddCartItemAutoID(c.Context(), cartID, it.ProductID, it.Quantity, price); err != nil {
			log.Printf("AddCartItemAutoID error: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "failed to add cart items"})
		}
	}

	// 3) Response must match doc exactly :contentReference[oaicite:7]{index=7}
	return c.Status(201).JSON(AddCartResponse{
		CartID:  cartID,
		Message: "Cart created and items added successfully",
	})
}

// GET /api/Carts/customer/:customer_id
// Response: []CartResponse :contentReference[oaicite:8]{index=8}
func GetCartByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customer_id")
	if customerID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "customer_id is required"})
	}

	carts, err := repo.GetCartsByCustomer(customerID)
	if err != nil {
		log.Printf("GetCartsByCustomer error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch carts"})
	}

	resp := make([]CartResponse, 0, len(carts))
	for _, cart := range carts {
		items, err := repo.GetCartItems(cart.CartID)
		if err != nil {
			log.Printf("GetCartItems error: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "failed to fetch cart items"})
		}

		itemResp := make([]CartItemResponse, 0, len(items))
		for _, it := range items {
			qty := 0
			if it.Quantity != nil {
				qty = *it.Quantity
			}
			price := 0.0
			if it.Price != nil {
				price = *it.Price
			}
			pid := ""
			if it.ProductID != nil {
				pid = *it.ProductID
			}

			itemResp = append(itemResp, CartItemResponse{
				CartItemID: it.CartItemID,
				ProductID:  pid,
				Quantity:   qty,
				Price:      price,
			})
		}

		cid := ""
		if cart.CustomerID != nil {
			cid = strings.TrimSpace(*cart.CustomerID)
		}
		st := ""
		if cart.Status != nil {
			st = strings.TrimSpace(*cart.Status)
		}

		resp = append(resp, CartResponse{
			CartID:     cart.CartID,
			CustomerID: cid,
			Status:     st,
			Items:      itemResp,
		})
	}

	return c.JSON(resp)
}

// GET /api/Carts/:cart_id
// Response: CartResponse :contentReference[oaicite:9]{index=9}
func GetCartByID(c *fiber.Ctx) error {
	cartID := c.Params("cart_id")
	if cartID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "cart_id is required"})
	}

	cart, err := repo.GetCartByID(cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": "cart not found"})
		}
		log.Printf("GetCartByID error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch cart"})
	}

	items, err := repo.GetCartItems(cartID)
	if err != nil {
		log.Printf("GetCartItems error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch cart items"})
	}

	itemResp := make([]CartItemResponse, 0, len(items))
	for _, it := range items {
		qty := 0
		if it.Quantity != nil {
			qty = *it.Quantity
		}
		price := 0.0
		if it.Price != nil {
			price = *it.Price
		}
		pid := ""
		if it.ProductID != nil {
			pid = *it.ProductID
		}

		itemResp = append(itemResp, CartItemResponse{
			CartItemID: it.CartItemID,
			ProductID:  pid,
			Quantity:   qty,
			Price:      price,
		})
	}

	cid := ""
	if cart.CustomerID != nil {
		cid = *cart.CustomerID
	}
	st := ""
	if cart.Status != nil {
		st = *cart.Status
	}

	return c.JSON(CartResponse{
		CartID:     cart.CartID,
		CustomerID: cid,
		Status:     st,
		Items:      itemResp,
	})
}

type UpdateCartStatusRequest struct {
	Status string `json:"status"`
}

func UpdateCartStatus(c *fiber.Ctx) error {
	cartID := c.Params("cart_id")

	var req UpdateCartStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}
	if req.Status == "" {
		return c.Status(400).JSON(fiber.Map{"message": "status is required"})
	}

	if err := repo.UpdateCartStatus(cartID, req.Status); err != nil {
		log.Printf("UpdateCartStatus error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to update cart status"})
	}

	return c.JSON(fiber.Map{
		"message": "cart status updated",
		"cart_id": cartID,
		"status":  req.Status,
	})
}
