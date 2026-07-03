package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailStatus string

const (
	StatusQueued  EmailStatus = "queued"
	StatusSent    EmailStatus = "sent"
	StatusFailed  EmailStatus = "failed"
)

type EmailLog struct {
	ID           uuid.UUID
	ApiKeyID     *uuid.UUID
	FromAddress  string
	ToAddresses  []string
	CcAddresses  []string
	BccAddresses []string
	Subject      string
	TemplateName *string
	Status       EmailStatus
	ErrorMsg     *string
	Attempts     int
	SentAt       *time.Time
	CreatedAt    time.Time
}
