package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, key *domain.ApiKey) error
	FindByHash(ctx context.Context, keyHash string) (*domain.ApiKey, error)
	FindAll(ctx context.Context) ([]domain.ApiKey, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type TemplateRepository interface {
	Create(ctx context.Context, t *domain.Template) error
	FindByName(ctx context.Context, name string) (*domain.Template, error)
	FindAll(ctx context.Context) ([]domain.Template, error)
	Update(ctx context.Context, t *domain.Template) error
	Delete(ctx context.Context, name string) error
}

type EmailLogRepository interface {
	Create(ctx context.Context, log *domain.EmailLog) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.EmailStatus, errMsg *string, attempts int) error
	MarkSent(ctx context.Context, id uuid.UUID, attempts int) error
	FindAll(ctx context.Context, filter dto.LogFilter) ([]domain.EmailLog, int, error)
}
