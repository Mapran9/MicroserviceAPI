package models

type CartItemDTO struct {
	CartItemID string  `json:"cart_item_id"`
	ProductID  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

type CartDTO struct {
	CartID     string        `json:"cart_id"`
	CustomerID string        `json:"customer_id"`
	Status     string        `json:"status"`
	Items      []CartItemDTO `json:"items"`
}
