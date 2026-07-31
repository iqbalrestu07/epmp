package entity

import "time"

// RefreshToken is used to issue new access tokens without re-authentication.
type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// IsValid returns true if the refresh token has not been revoked and has not expired.
func (r *RefreshToken) IsValid() bool {
	return r.RevokedAt == nil && r.ExpiresAt.After(time.Now())
}
