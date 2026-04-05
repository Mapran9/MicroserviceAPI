package repo

import (
	"context"
	"customer/internal/models"
	"database/sql"
)

const (
	customerPrefix   = "CTM"
	customerDigits   = 9
	customerLockName = "customers_id_lock"
)

func GetAllCustomers() ([]models.Customer, error) {
	rows, err := DB.Query(`
		SELECT
			customer_id,
			first_name,
			last_name,
			email,
			phone_number,
			address,
			city,
			state,
			country,
			postal_code,
			created_at
		FROM customers
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(
			&c.CustomerID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.PhoneNumber,
			&c.Address,
			&c.City,
			&c.State,
			&c.Country,
			&c.PostalCode,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, nil
}

func GetCustomerByID(id string) (*models.Customer, error) {
	row := DB.QueryRow(`
		SELECT
			customer_id,
			first_name,
			last_name,
			email,
			phone_number,
			address,
			city,
			state,
			country,
			postal_code,
			created_at
		FROM customers
		WHERE customer_id = ?
	`, id)

	var c models.Customer
	if err := row.Scan(
		&c.CustomerID,
		&c.FirstName,
		&c.LastName,
		&c.Email,
		&c.PhoneNumber,
		&c.Address,
		&c.City,
		&c.State,
		&c.Country,
		&c.PostalCode,
		&c.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &c, nil
}

// Auto ID version (uses id_lock_conn.go)
func CreateCustomerAutoID(ctx context.Context, c *models.Customer) (string, error) {
	var newID string

	err := WithAdvisoryLockConn(ctx, DB, customerLockName, func(conn *sql.Conn) error {
		// 1) read last id (เห็นเฉพาะ committed rows)
		var lastID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT customer_id
			FROM customers
			WHERE customer_id LIKE CONCAT(?, '%')
			ORDER BY customer_id DESC
			LIMIT 1
		`, customerPrefix).Scan(&lastID)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		// 2) next id
		newID = NextIDFromLast(lastID, customerPrefix, customerDigits)

		// 3) insert (autocommit => commit ทันที ก่อนปล่อย lock)
		_, err = conn.ExecContext(ctx, `
			INSERT INTO customers (
				customer_id,
				first_name,
				last_name,
				email,
				phone_number,
				address,
				city,
				state,
				country,
				postal_code
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			newID,
			c.FirstName,
			c.LastName,
			c.Email,
			c.PhoneNumber,
			c.Address,
			c.City,
			c.State,
			c.Country,
			c.PostalCode,
		)
		return err
	})

	return newID, err
}
