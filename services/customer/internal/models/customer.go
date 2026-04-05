package models

import "time"

type Customer struct {
	CustomerID  string    `json:"customer_id"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	Email       *string   `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	Address     *string   `json:"address"`
	City        *string   `json:"city"`
	State       *string   `json:"state"`
	Country     *string   `json:"country"`
	PostalCode  *string   `json:"postal_code"`
	CreatedAt   time.Time `json:"created_at"`
}
