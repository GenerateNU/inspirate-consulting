package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"inspirate-consulting/internal/errs"

	"github.com/google/uuid"
)

// TimeIDKey is the sort key for tables paginated by (created_at, id).
type TimeIDKey struct {
	CreatedAt time.Time `json:"t"`
	ID        uuid.UUID `json:"i"`
}

func Encode[K any](key K) (string, error) {
	b, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Decode[K any](cursor string) (K, error) {
	var key K
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return key, errs.BadRequest("invalid cursor")
	}
	if err := json.Unmarshal(b, &key); err != nil {
		return key, errs.BadRequest("invalid cursor")
	}
	return key, nil
}

// NextPage expects rows fetched with limit+1; the extra row only signals that another page exists.
func NextPage[T any, K any](rows []T, limit int, key func(T) K) ([]T, *string, error) {
	if len(rows) <= limit {
		return rows, nil, nil
	}
	page := rows[:limit]
	cursor, err := Encode(key(page[limit-1]))
	if err != nil {
		return nil, nil, err
	}
	return page, &cursor, nil
}
