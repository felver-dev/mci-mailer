package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcicare/mci-mailer/internal/domain"
)

type apiKeyRepo struct {
	db *pgxpool.Pool
}

func NewApiKeyRepository(db *pgxpool.Pool) ApiKeyRepository {
	return &apiKeyRepo{db: db}
}

func (r *apiKeyRepo) Create(ctx context.Context, key *domain.ApiKey) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO api_keys (id, name, key_hash, scopes, is_active, created_at, created_by_user_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		key.ID, key.Name, key.KeyHash, key.Scopes, key.IsActive, key.CreatedAt, key.CreatedByUserID,
	)
	return err
}

func (r *apiKeyRepo) FindByHash(ctx context.Context, keyHash string) (*domain.ApiKey, error) {
	k := &domain.ApiKey{}
	err := r.db.QueryRow(ctx,
		`SELECT id, name, key_hash, scopes, is_active, created_at, last_used_at, created_by_user_id
		 FROM api_keys WHERE key_hash = $1 AND is_active = TRUE`,
		keyHash,
	).Scan(&k.ID, &k.Name, &k.KeyHash, &k.Scopes, &k.IsActive, &k.CreatedAt, &k.LastUsedAt, &k.CreatedByUserID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return k, err
}

func (r *apiKeyRepo) FindAll(ctx context.Context) ([]domain.ApiKey, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ak.id, ak.name, ak.key_hash, ak.scopes, ak.is_active, ak.created_at, ak.last_used_at,
		        ak.created_by_user_id, u.name AS created_by_name
		 FROM api_keys ak
		 LEFT JOIN users u ON u.id = ak.created_by_user_id
		 ORDER BY ak.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []domain.ApiKey
	for rows.Next() {
		var k domain.ApiKey
		if err := rows.Scan(
			&k.ID, &k.Name, &k.KeyHash, &k.Scopes, &k.IsActive, &k.CreatedAt, &k.LastUsedAt,
			&k.CreatedByUserID, &k.CreatedByUserName,
		); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *apiKeyRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE api_keys SET is_active = FALSE WHERE id = $1`, id,
	)
	return err
}

func (r *apiKeyRepo) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE api_keys SET last_used_at = NOW() WHERE id = $1`, id,
	)
	return err
}
