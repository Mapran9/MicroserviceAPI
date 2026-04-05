package models

import "time"

type OrderDetail struct {
	OrderItemID string    `json:"order_item_id"`
	OrderID     *string   `json:"order_id"`
	ProductID   *string   `json:"product_id"`
	Quantity    *int      `json:"quantity"`
	UnitPrice   *float64  `json:"unit_price"`
	TotalPrice  *float64  `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}
