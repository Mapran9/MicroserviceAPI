package repo

import (
	"context"
	"database/sql"

	"cart/internal/models"
)

const (
	cartPrefix   = "CART"
	cartDigits   = 6
	cartLockName = "carts_id_lock"

	cartItemPrefix   = "CI"
	cartItemDigits   = 8
	cartItemLockName = "cart_items_id_lock"
)

// Create cart header (CART000001...) with provided status
func CreateCartAutoID(ctx context.Context, customerID string, status string) (string, error) {
	var newID string

	err := WithAdvisoryLockConn(ctx, DB, cartLockName, func(conn *sql.Conn) error {
		var lastID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT cart_id
			FROM carts
			WHERE cart_id LIKE CONCAT(?, '%')
			ORDER BY cart_id DESC
			LIMIT 1
		`, cartPrefix).Scan(&lastID)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		newID = NextIDFromLast(lastID, cartPrefix, cartDigits)

		_, err = conn.ExecContext(ctx, `
			INSERT INTO carts (cart_id, customer_id, status)
			VALUES (?, ?, ?)
		`, newID, customerID, status)

		return err
	})

	return newID, err
}

// Add cart item (CI00000001...)
func AddCartItemAutoID(ctx context.Context, cartID string, productID string, qty int, price float64) (string, error) {
	var newID string

	err := WithAdvisoryLockConn(ctx, DB, cartItemLockName, func(conn *sql.Conn) error {
		var lastID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT cart_item_id
			FROM cart_items
			WHERE cart_item_id LIKE CONCAT(?, '%')
			ORDER BY cart_item_id DESC
			LIMIT 1
		`, cartItemPrefix).Scan(&lastID)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		newID = NextIDFromLast(lastID, cartItemPrefix, cartItemDigits)

		_, err = conn.ExecContext(ctx, `
			INSERT INTO cart_items (cart_item_id, cart_id, product_id, quantity, price)
			VALUES (?, ?, ?, ?, ?)
		`, newID, cartID, productID, qty, price)

		return err
	})

	return newID, err
}

func GetCartByID(cartID string) (*models.Cart, error) {
	row := DB.QueryRow(`
		SELECT cart_id, customer_id, status, created_at
		FROM carts
		WHERE cart_id = ?
	`, cartID)

	var c models.Cart
	if err := row.Scan(&c.CartID, &c.CustomerID, &c.Status, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func GetCartsByCustomer(customerID string) ([]models.Cart, error) {
	rows, err := DB.Query(`
		SELECT cart_id, customer_id, status, created_at
		FROM carts
		WHERE customer_id = ?
		ORDER BY created_at DESC
	`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	carts := []models.Cart{}
	for rows.Next() {
		var c models.Cart
		if err := rows.Scan(&c.CartID, &c.CustomerID, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		carts = append(carts, c)
	}
	return carts, nil
}

func GetCartItems(cartID string) ([]models.CartItem, error) {
	rows, err := DB.Query(`
		SELECT cart_item_id, cart_id, product_id, quantity, price, created_at
		FROM cart_items
		WHERE cart_id = ?
		ORDER BY created_at DESC
	`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.CartItemID, &it.CartID, &it.ProductID, &it.Quantity, &it.Price, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}
func UpdateCartStatus(cartID string, status string) error {
	_, err := DB.Exec(`UPDATE carts SET status = ? WHERE cart_id = ?`, status, cartID)
	return err
}
