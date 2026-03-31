package sms

import (
	"context"
	"log"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SMSRepository interface {
	EnqueueSMS(ctx context.Context, msisdn, rawMessage string) error
	GetPendingSMSBatch(ctx context.Context, limit int) ([]models.SMSQueueItem, error)
	MarkSMSAsProcessed(ctx context.Context, id uuid.UUID) error
	MarkSMSRejected(ctx context.Context, id uuid.UUID, reason string) error
}

var SMSRepo SMSRepository

func InitSMSRepo(db *pgxpool.Pool, logger *log.Logger) {
	SMSRepo = newMgRepository(db, logger)
}
