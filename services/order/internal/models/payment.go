package models

import "time"

type Payment struct {
	PaymentID     string    `json:"payment_id"`
	CustomerID    *string   `json:"customer_id"`
	OrderID       *string   `json:"order_id"`
	PaymentMethod *string   `json:"payment_method"`
	Amount        *float64  `json:"amount"`
	Status        *string   `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
