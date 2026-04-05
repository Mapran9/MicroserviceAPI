package repo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

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
