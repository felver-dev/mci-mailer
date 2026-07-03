package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID          uuid.UUID
	Name        string
	KeyHash     string
	Scopes      []string
	IsActive    bool
	CreatedAt   time.Time
	LastUsedAt  *time.Time
}

const (
	ScopeMailSend       = "mail:send"
	ScopeTemplatesRead  = "templates:read"
	ScopeTemplatesWrite = "templates:write"
	ScopeLogsRead       = "logs:read"
	ScopeKeysManage     = "keys:manage"
)

func GenerateAPIKey() (rawKey, keyHash string, err error) {
	b := make([]byte, 24)
	if _, err = rand.Read(b); err != nil {
		return
	}
	rawKey = "MCM." + hex.EncodeToString(b)
	keyHash = HashAPIKey(rawKey)
	return
}

func HashAPIKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(h[:])
}

func (k *ApiKey) HasScope(scope string) bool {
	for _, s := range k.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}
