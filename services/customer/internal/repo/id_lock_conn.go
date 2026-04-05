package repo

import (
	"context"
	"database/sql"
	"fmt"
)

func WithAdvisoryLockConn(ctx context.Context, db *sql.DB, lockName string, fn func(conn *sql.Conn) error) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var got int
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, 5)", lockName).Scan(&got); err != nil {
		return err
	}
	if got != 1 {
		return fmt.Errorf("could not acquire lock: %s", lockName)
	}

	// ปล่อย lock หลังงานเสร็จ (ยังอยู่บน conn เดิมแน่นอน)
	defer func() { _, _ = conn.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", lockName) }()

	return fn(conn)
}
