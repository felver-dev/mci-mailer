package domain

import (
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID        uuid.UUID
	Name      string
	Subject   string
	HtmlBody  string
	TextBody  string
	Variables []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
