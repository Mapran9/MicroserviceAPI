package models

import "time"

type Cart struct {
	CartID     string     `json:"cart_id"`
	CustomerID *string    `json:"customer_id"`
	Status     *string    `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	Items      []CartItem `json:"items,omitempty"`
}
