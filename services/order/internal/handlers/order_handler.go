package handlers

import (
	"database/sql"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"order/internal/clients"
	"order/internal/repo"
)

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	CartID     string `json:"cart_id"`
	Payment    struct {
		PaymentMethod string `json:"payment_method"`
	} `json:"payment"`
}

type CreateOrderResponse struct {
	Message    string  `json:"message"`
	OrderID    string  `json:"order_id"`
	OrderTotal float64 `json:"order_total"`
	PaymentID  string  `json:"payment_id"`
	Status     string  `json:"status"`
}

type OrderItemResponse struct {
	OrderItemID string  `json:"order_item_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderResponse struct {
	OrderID     string              `json:"order_id"`
	CustomerID  string              `json:"customer_id"`
	CartID      string              `json:"cart_id"`
	OrderTotal  float64             `json:"order_total"`
	OrderStatus string              `json:"order_status"`
	Items       []OrderItemResponse `json:"items,omitempty"`
}

func CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	if req.CustomerID == "" || req.CartID == "" || req.Payment.PaymentMethod == "" {
		return c.Status(400).JSON(fiber.Map{"message": "customer_id, cart_id and payment.payment_method are required"})
	}

	cart, err := clients.GetCart(req.CartID)
	if err != nil {
		log.Printf("GetCart error: %v", err)
		return c.Status(400).JSON(fiber.Map{"message": "invalid cart_id"})
	}

	if strings.TrimSpace(cart.CustomerID) != strings.TrimSpace(req.CustomerID) {
		return c.Status(400).JSON(fiber.Map{"message": "customer_id does not match cart owner"})
	}

	if len(cart.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "cart has no items"})
	}

	orderID, total, err := repo.CreateOrderFromCart(c.Context(), *cart)
	if err != nil {
		log.Printf("CreateOrderFromCart error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to create order"})
	}

	paymentResp, err := clients.CreatePayment(clients.CreatePaymentRequest{
		CustomerID:    req.CustomerID,
		OrderID:       orderID,
		PaymentMethod: req.Payment.PaymentMethod,
		Amount:        total,
		Status:        "pending",
	})
	if err != nil {
		log.Printf("CreatePayment error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to create payment"})
	}

	if err := clients.UpdateCartStatus(req.CartID, "ordered"); err != nil {
		log.Printf("UpdateCartStatus error: %v", err)
	}

	return c.Status(201).JSON(CreateOrderResponse{
		Message:    "Order created successfully",
		OrderID:    orderID,
		OrderTotal: total,
		PaymentID:  paymentResp.PaymentID,
		Status:     "pending",
	})
}

func GetOrderByID(c *fiber.Ctx) error {
	orderID := c.Params("order_id")

	order, err := repo.GetOrderByID(orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"message": "order not found"})
		}
		log.Printf("GetOrderByID error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch order"})
	}

	items, err := repo.GetOrderItems(orderID)
	if err != nil {
		log.Printf("GetOrderItems error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch order items"})
	}

	respItems := make([]OrderItemResponse, 0, len(items))
	for _, it := range items {
		pid := ""
		if it.ProductID != nil {
			pid = strings.TrimSpace(*it.ProductID)
		}
		qty := 0
		if it.Quantity != nil {
			qty = *it.Quantity
		}
		up := 0.0
		if it.UnitPrice != nil {
			up = *it.UnitPrice
		}
		tp := 0.0
		if it.TotalPrice != nil {
			tp = *it.TotalPrice
		}

		respItems = append(respItems, OrderItemResponse{
			OrderItemID: it.OrderItemID,
			ProductID:   pid,
			Quantity:    qty,
			UnitPrice:   up,
			TotalPrice:  tp,
		})
	}

	customerID := ""
	if order.CustomerID != nil {
		customerID = strings.TrimSpace(*order.CustomerID)
	}
	cartID := ""
	if order.CartID != nil {
		cartID = strings.TrimSpace(*order.CartID)
	}
	status := ""
	if order.OrderStatus != nil {
		status = strings.TrimSpace(*order.OrderStatus)
	}
	total := 0.0
	if order.OrderTotal != nil {
		total = *order.OrderTotal
	}

	return c.JSON(OrderResponse{
		OrderID:     order.OrderID,
		CustomerID:  customerID,
		CartID:      cartID,
		OrderTotal:  total,
		OrderStatus: status,
		Items:       respItems,
	})
}

func GetOrdersByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customer_id")

	orders, err := repo.GetOrdersByCustomer(customerID)
	if err != nil {
		log.Printf("GetOrdersByCustomer error: %v", err)
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch orders"})
	}

	resp := make([]OrderResponse, 0, len(orders))

	for _, order := range orders {
		cid := ""
		if order.CustomerID != nil {
			cid = strings.TrimSpace(*order.CustomerID)
		}
		cartID := ""
		if order.CartID != nil {
			cartID = strings.TrimSpace(*order.CartID)
		}
		status := ""
		if order.OrderStatus != nil {
			status = strings.TrimSpace(*order.OrderStatus)
		}
		total := 0.0
		if order.OrderTotal != nil {
			total = *order.OrderTotal
		}

		items, err := repo.GetOrderItems(order.OrderID)
		if err != nil {
			log.Printf("GetOrderItems error: %v", err)
			return c.Status(500).JSON(fiber.Map{"message": "failed to fetch order items"})
		}

		respItems := make([]OrderItemResponse, 0, len(items))
		for _, it := range items {
			pid := ""
			if it.ProductID != nil {
				pid = strings.TrimSpace(*it.ProductID)
			}
			qty := 0
			if it.Quantity != nil {
				qty = *it.Quantity
			}
			up := 0.0
			if it.UnitPrice != nil {
				up = *it.UnitPrice
			}
			tp := 0.0
			if it.TotalPrice != nil {
				tp = *it.TotalPrice
			}

			respItems = append(respItems, OrderItemResponse{
				OrderItemID: it.OrderItemID,
				ProductID:   pid,
				Quantity:    qty,
				UnitPrice:   up,
				TotalPrice:  tp,
			})
		}

		resp = append(resp, OrderResponse{
			OrderID:     order.OrderID,
			CustomerID:  cid,
			CartID:      cartID,
			OrderTotal:  total,
			OrderStatus: status,
			Items:       respItems,
		})
	}

	return c.JSON(resp)
}
