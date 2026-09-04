package uid

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// New generates a new ULID as an uppercase string (26 chars).
// ULIDs are lexicographically sortable and time-ordered.
// Stored in PostgreSQL as TEXT.
func New() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
}

// IsValid checks whether a string is a valid ULID.
func IsValid(s string) bool {
	_, err := ulid.ParseStrict(s)
	return err == nil
}
