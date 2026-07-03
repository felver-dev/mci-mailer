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

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindAll(ctx context.Context) ([]domain.User, error)
	Count(ctx context.Context) (int, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

type StatsRepository interface {
	GetOverview(ctx context.Context) (*dto.StatsOverview, error)
	GetPerApp(ctx context.Context) ([]dto.AppStats, error)
}
