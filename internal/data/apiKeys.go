package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type APIKey struct {
	ID         int64      `json:"id"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	Name       string     `json:"name"`
	KeyHash    string     `json:"-"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type APIKeyModel struct {
	DB *sql.DB
}

func (m APIKeyModel) GetByHash(hash string) (*APIKey, error) {
	query := `
        SELECT id, name, is_active, expires_at
        FROM api_keys
        WHERE key_hash = $1 AND is_active = true
    `

	var apiKey APIKey
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, hash).Scan(
		&apiKey.ID,
		&apiKey.Name,
		&apiKey.IsActive,
		&apiKey.ExpiresAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if apiKey.ExpiresAt != nil && time.Now().After(*apiKey.ExpiresAt) {
		return nil, ErrRecordNotFound
	}

	return &apiKey, nil
}

func (m APIKeyModel) UpdateLastUsed(id int64) error {
	query := `UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
