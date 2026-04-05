package models

import "time"

type CartItem struct {
	CartItemID string    `json:"cart_item_id"`
	CartID     *string   `json:"cart_id"`
	ProductID  *string   `json:"product_id"`
	Quantity   *int      `json:"quantity"`
	Price      *float64  `json:"price"` // DECIMAL scan -> float64 ok
	CreatedAt  time.Time `json:"created_at"`
}
