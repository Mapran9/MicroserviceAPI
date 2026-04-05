package repo

import (
	"context"
	"database/sql"
)

const (
	paymentPrefix   = "PM"
	paymentDigits   = 5
	paymentLockName = "payments_id_lock"
)

func CreatePaymentAutoID(ctx context.Context, customerID, orderID, paymentMethod string, amount float64, status string) (string, error) {
	var newPaymentID string

	err := WithAdvisoryLockConn(ctx, DB, paymentLockName, func(conn *sql.Conn) error {
		var lastPaymentID sql.NullString
		err := conn.QueryRowContext(ctx, `
			SELECT payment_id
			FROM payments
			WHERE payment_id LIKE CONCAT(?, '%')
			ORDER BY payment_id DESC
			LIMIT 1
		`, paymentPrefix).Scan(&lastPaymentID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		newPaymentID = NextIDFromLast(lastPaymentID, paymentPrefix, paymentDigits)

		_, err = conn.ExecContext(ctx, `
			INSERT INTO payments (
				payment_id, customer_id, order_id, payment_method, amount, status
			) VALUES (?, ?, ?, ?, ?, ?)
		`, newPaymentID, customerID, orderID, paymentMethod, amount, status)

		return err
	})

	return newPaymentID, err
}
