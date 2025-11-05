package mfa

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOTPStore struct {
	db *pgxpool.Pool
}

// NewPostgresOTPStore creates a new Postgres OTP store
func NewPostgresOTPStore(db *pgxpool.Pool) *PostgresOTPStore {
	return &PostgresOTPStore{db: db}
}

// Set stores an OTP with an expiry
func (p *PostgresOTPStore) Set(key, code string, ttl time.Duration) error {
	ctx := context.Background()
	expiry := time.Now().Add(ttl)

	query := `
		INSERT INTO otps (key, code, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE
		SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at
	`

	_, err := p.db.Exec(ctx, query, key, code, expiry)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}
	return nil
}

// Get retrieves an OTP if it's not expired
func (p *PostgresOTPStore) Get(key string) (string, bool) {
	ctx := context.Background()

	query := `
		SELECT code, expires_at FROM otps WHERE key = $1
	`

	var code string
	var expiresAt time.Time
	err := p.db.QueryRow(ctx, query, key).Scan(&code, &expiresAt)
	if err != nil {
		return "", false
	}

	if time.Now().After(expiresAt) {
		p.Delete(key)
		return "", false
	}

	return code, true
}

// Delete removes an OTP entry
func (p *PostgresOTPStore) Delete(key string) error {
	ctx := context.Background()
	_, err := p.db.Exec(ctx, `DELETE FROM otps WHERE key = $1`, key)
	return err
}

// CleanupExpired removes old OTPs (optional)
func (p *PostgresOTPStore) CleanupExpired() error {
	ctx := context.Background()
	_, err := p.db.Exec(ctx, `DELETE FROM otps WHERE expires_at < NOW()`)
	return err
}
