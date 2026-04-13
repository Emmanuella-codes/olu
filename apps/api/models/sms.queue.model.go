package models

import "github.com/google/uuid"

type SMSQueueItem struct {
	ID         uuid.UUID
	MSISDN     string
	RawMessage string
}
