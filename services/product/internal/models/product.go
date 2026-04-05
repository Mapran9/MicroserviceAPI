package models

import "time"

type Product struct {
	ProductID   string    `json:"product_id"`
	ProductName *string   `json:"product_name"`
	Brand       *string   `json:"brand"`
	Category    *string   `json:"category"`
	Price       *float64  `json:"price"`
	Stock       *int      `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}
