package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mcicare/mci-mailer/internal/domain"
	"github.com/mcicare/mci-mailer/internal/dto"
	"github.com/mcicare/mci-mailer/internal/repository"
)

const tokenTTL = 24 * time.Hour

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type UserService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewUserService(repo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *UserService) Bootstrap(ctx context.Context, req dto.BootstrapRequest) (*dto.LoginResponse, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("portal already initialized — use login")
	}
	return s.createUser(ctx, req.Name, req.Email, req.Password, domain.RoleAdmin)
}

func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := domain.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         domain.UserRole(req.Role),
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	resp := toUserResponse(u)
	return &resp, nil
}

func (s *UserService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil || !domain.CheckPassword(u.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	go func() { _ = s.repo.UpdateLastLogin(context.Background(), u.ID) }()

	return s.issueToken(u)
}

func (s *UserService) List(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.UserResponse, len(users))
	for i, u := range users {
		result[i] = toUserResponse(&u)
	}
	return result, nil
}

func (s *UserService) ValidateToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ── private helpers ──────────────────────────────────────────────────────────

func (s *UserService) createUser(ctx context.Context, name, email, password string, role domain.UserRole) (*dto.LoginResponse, error) {
	hash, err := domain.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return s.issueToken(u)
}

func (s *UserService) issueToken(u *domain.User) (*dto.LoginResponse, error) {
	exp := time.Now().Add(tokenTTL)
	claims := JWTClaims{
		UserID: u.ID.String(),
		Email:  u.Email,
		Role:   string(u.Role),
		Name:   u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}
	resp := toUserResponse(u)
	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: exp.UTC().Format(time.RFC3339),
		User:      resp,
	}, nil
}

func toUserResponse(u *domain.User) dto.UserResponse {
	r := dto.UserResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	}
	if u.LastLoginAt != nil {
		s := u.LastLoginAt.UTC().Format(time.RFC3339)
		r.LastLoginAt = &s
	}
	return r
}
