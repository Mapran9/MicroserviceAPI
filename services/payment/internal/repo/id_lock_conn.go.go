package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// WithAdvisoryLockTx runs fn(tx) inside:
// - BEGIN
// - GET_LOCK(lockName)
// - fn(tx) [generate + insert]
// - RELEASE_LOCK(lockName)
// - COMMIT
// If error -> ROLLBACK (lock is released when connection closes; we also try RELEASE_LOCK)
func WithAdvisoryLockTx(ctx context.Context, db *sql.DB, lockName string, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// acquire lock (wait up to 3 seconds)
	var got int
	if err := tx.QueryRowContext(ctx, "SELECT GET_LOCK(?, 3)", lockName).Scan(&got); err != nil {
		return err
	}
	if got != 1 {
		return fmt.Errorf("could not acquire lock: %s", lockName)
	}

	// ensure unlock attempt (same connection/tx)
	defer func() {
		_, _ = tx.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", lockName)
	}()

	// run critical section
	if err := fn(tx); err != nil {
		return err
	}

	// release lock before commit (ok either way; keep it explicit)
	if _, err := tx.ExecContext(ctx, "SELECT RELEASE_LOCK(?)", lockName); err != nil {
		return err
	}

	return tx.Commit()
}

// NextIDFromLast generates prefix + zero padded digits.
// lastID must be fixed format: prefix + digits, e.g. CTM000000001
func NextIDFromLast(lastID sql.NullString, prefix string, digits int) string {
	nextNum := int64(1)
	if lastID.Valid && strings.HasPrefix(lastID.String, prefix) {
		numPart := lastID.String[len(prefix):]
		if len(numPart) == digits {
			if n, err := strconv.ParseInt(numPart, 10, 64); err == nil {
				nextNum = n + 1
			}
		}
	}
	return fmt.Sprintf("%s%0*d", prefix, digits, nextNum)
}
