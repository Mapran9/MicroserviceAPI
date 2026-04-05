package repo

import (
	"context"
	"database/sql"

	"product/internal/models"
)

const (
	productPrefix   = "P"
	productDigits   = 5
	productLockName = "products_id_lock"
)

func GetAllProducts() ([]models.Product, error) {
	rows, err := DB.Query(`
		SELECT
			product_id,
			product_name,
			brand,
			category,
			price,
			stock,
			created_at
		FROM products
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ProductID,
			&p.ProductName,
			&p.Brand,
			&p.Category,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(id string) (*models.Product, error) {
	row := DB.QueryRow(`
		SELECT
			product_id,
			product_name,
			brand,
			category,
			price,
			stock,
			created_at
		FROM products
		WHERE product_id = ?
	`, id)

	var p models.Product
	if err := row.Scan(
		&p.ProductID,
		&p.ProductName,
		&p.Brand,
		&p.Category,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

// Auto ID + concurrent-safe (GET_LOCK + same conn + SELECT last + INSERT)
func CreateProductAutoID(ctx context.Context, p *models.Product) (string, error) {
	var newID string

	err := WithAdvisoryLockConn(ctx, DB, productLockName, func(conn *sql.Conn) error {
		var lastID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT product_id
			FROM products
			WHERE product_id LIKE CONCAT(?, '%')
			ORDER BY product_id DESC
			LIMIT 1
		`, productPrefix).Scan(&lastID)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		newID = NextIDFromLast(lastID, productPrefix, productDigits)

		_, err = conn.ExecContext(ctx, `
			INSERT INTO products (
				product_id,
				product_name,
				brand,
				category,
				price,
				stock
			) VALUES (?, ?, ?, ?, ?, ?)
		`,
			newID,
			p.ProductName,
			p.Brand,
			p.Category,
			p.Price,
			p.Stock,
		)
		return err
	})

	return newID, err
}
