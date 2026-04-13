package models

import (
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsActive     bool
	LastLogin    *time.Time
}
