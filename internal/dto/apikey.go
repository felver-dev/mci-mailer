package dto

type CreateApiKeyRequest struct {
	Name   string   `json:"name" binding:"required"`
	Scopes []string `json:"scopes" binding:"required,min=1"`
}

type CreateApiKeyResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Key       string   `json:"key"` // affiché une seule fois
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"created_at"`
}

type ApiKeyResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Scopes      []string `json:"scopes"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
	LastUsedAt  *string  `json:"last_used_at"`
}
