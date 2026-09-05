package uid

import (
	"github.com/google/uuid"
)

// New generates a new time-ordered UUIDv7 string (36 chars with hyphens).
// UUIDv7 is compliant with RFC 9562, combines 48-bit UNIX timestamp (ordered)
// and cryptographically random data, stored in PostgreSQL as native UUID.
func New() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}

// IsValid checks whether a string is a valid UUID.
func IsValid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
