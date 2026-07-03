package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

type ApiKeyService struct {
	repo repository.ApiKeyRepository
}

func NewApiKeyService(repo repository.ApiKeyRepository) *ApiKeyService {
	return &ApiKeyService{repo: repo}
}

func (s *ApiKeyService) Create(ctx context.Context, req dto.CreateApiKeyRequest) (*dto.CreateApiKeyResponse, error) {
	if err := validateScopes(req.Scopes); err != nil {
		return nil, err
	}

	rawKey, keyHash, err := domain.GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	key := &domain.ApiKey{
		ID:        uuid.New(),
		Name:      req.Name,
		KeyHash:   keyHash,
		Scopes:    req.Scopes,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, key); err != nil {
		return nil, err
	}

	return &dto.CreateApiKeyResponse{
		ID:        key.ID.String(),
		Name:      key.Name,
		Key:       rawKey,
		Scopes:    key.Scopes,
		CreatedAt: key.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *ApiKeyService) List(ctx context.Context) ([]dto.ApiKeyResponse, error) {
	keys, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ApiKeyResponse, len(keys))
	for i, k := range keys {
		result[i] = dto.ApiKeyResponse{
			ID:        k.ID.String(),
			Name:      k.Name,
			Scopes:    k.Scopes,
			IsActive:  k.IsActive,
			CreatedAt: k.CreatedAt.Format(time.RFC3339),
		}
		if k.LastUsedAt != nil {
			t := k.LastUsedAt.Format(time.RFC3339)
			result[i].LastUsedAt = &t
		}
	}
	return result, nil
}

func (s *ApiKeyService) Revoke(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid key id")
	}
	return s.repo.Revoke(ctx, uid)
}

func validateScopes(scopes []string) error {
	valid := map[string]bool{
		domain.ScopeMailSend:       true,
		domain.ScopeTemplatesRead:  true,
		domain.ScopeTemplatesWrite: true,
		domain.ScopeLogsRead:       true,
		domain.ScopeKeysManage:     true,
	}
	for _, s := range scopes {
		if !valid[s] {
			return errors.New("invalid scope: " + s)
		}
	}
	return nil
}
