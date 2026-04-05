package repo

import (
	"context"
	"database/sql"

	"order/internal/models"
)

const (
	orderPrefix   = "OR"
	orderDigits   = 8
	orderLockName = "orders_header_id_lock"

	orderItemPrefix = "OD"
	orderItemDigits = 8
)

func CreateOrderFromCart(ctx context.Context, cart models.CartDTO) (string, float64, error) {
	var newOrderID string
	var orderTotal float64

	err := WithAdvisoryLockConn(ctx, DB, orderLockName, func(conn *sql.Conn) error {
		var lastOrderID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT order_id
			FROM orders_header
			WHERE order_id LIKE CONCAT(?, '%')
			ORDER BY order_id DESC
			LIMIT 1
		`, orderPrefix).Scan(&lastOrderID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		newOrderID = NextIDFromLast(lastOrderID, orderPrefix, orderDigits)

		for _, item := range cart.Items {
			orderTotal += float64(item.Quantity) * item.Price
		}

		_, err = conn.ExecContext(ctx, `
			INSERT INTO orders_header (
				order_id, customer_id, cart_id, order_total, order_status
			) VALUES (?, ?, ?, ?, ?)
		`, newOrderID, cart.CustomerID, cart.CartID, orderTotal, "pending")
		if err != nil {
			return err
		}

		for _, item := range cart.Items {
			newOrderItemID, err := nextOrderItemID(ctx, conn)
			if err != nil {
				return err
			}

			totalPrice := float64(item.Quantity) * item.Price

			_, err = conn.ExecContext(ctx, `
				INSERT INTO order_detail (
					order_item_id, order_id, product_id, quantity, unit_price, total_price
				) VALUES (?, ?, ?, ?, ?, ?)
			`, newOrderItemID, newOrderID, item.ProductID, item.Quantity, item.Price, totalPrice)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return newOrderID, orderTotal, err
}

func nextOrderItemID(ctx context.Context, conn *sql.Conn) (string, error) {
	var lastID sql.NullString
	err := conn.QueryRowContext(ctx, `
		SELECT order_item_id
		FROM order_detail
		WHERE order_item_id LIKE CONCAT(?, '%')
		ORDER BY order_item_id DESC
		LIMIT 1
	`, orderItemPrefix).Scan(&lastID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return NextIDFromLast(lastID, orderItemPrefix, orderItemDigits), nil
}

func GetOrderByID(orderID string) (*models.OrderHeader, error) {
	row := DB.QueryRow(`
		SELECT order_id, customer_id, cart_id, order_total, order_status, created_at, updated_at
		FROM orders_header
		WHERE order_id = ?
	`, orderID)

	var o models.OrderHeader
	if err := row.Scan(
		&o.OrderID,
		&o.CustomerID,
		&o.CartID,
		&o.OrderTotal,
		&o.OrderStatus,
		&o.CreatedAt,
		&o.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &o, nil
}

func GetOrderItems(orderID string) ([]models.OrderDetail, error) {
	rows, err := DB.Query(`
		SELECT order_item_id, order_id, product_id, quantity, unit_price, total_price, created_at
		FROM order_detail
		WHERE order_id = ?
		ORDER BY created_at ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.OrderDetail{}
	for rows.Next() {
		var it models.OrderDetail
		if err := rows.Scan(
			&it.OrderItemID,
			&it.OrderID,
			&it.ProductID,
			&it.Quantity,
			&it.UnitPrice,
			&it.TotalPrice,
			&it.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}

	return items, nil
}

func GetOrdersByCustomer(customerID string) ([]models.OrderHeader, error) {
	rows, err := DB.Query(`
		SELECT order_id, customer_id, cart_id, order_total, order_status, created_at, updated_at
		FROM orders_header
		WHERE customer_id = ?
		ORDER BY created_at DESC
	`, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.OrderHeader{}
	for rows.Next() {
		var o models.OrderHeader
		if err := rows.Scan(
			&o.OrderID,
			&o.CustomerID,
			&o.CartID,
			&o.OrderTotal,
			&o.OrderStatus,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}
