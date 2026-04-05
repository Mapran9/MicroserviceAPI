package models

import "time"

type OrderHeader struct {
	OrderID     string    `json:"order_id"`
	CustomerID  *string   `json:"customer_id"`
	CartID      *string   `json:"cart_id"`
	OrderTotal  *float64  `json:"order_total"`
	OrderStatus *string   `json:"order_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
