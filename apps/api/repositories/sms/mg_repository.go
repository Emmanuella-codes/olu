package sms

import (
	"context"
	"log"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type mgRepository struct {
	db  *pgxpool.Pool
	log *log.Logger
}

func newMgRepository(db *pgxpool.Pool, logger *log.Logger) *mgRepository {
	return &mgRepository{
		db:  db,
		log: logger,
	}
}

func (r *mgRepository) EnqueueSMS(ctx context.Context, msisdn, rawMessage string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sms_queue (msisdn, raw_message, status, created_at)
		VALUES ($1, $2, 'pending', NOW())
	`, msisdn, rawMessage)
	return err
}

func (r *mgRepository) GetPendingSMSBatch(ctx context.Context, limit int) ([]models.SMSQueueItem, error) {
	rows, err := r.db.Query(ctx, `
		UPDATE sms_queue
		SET status = 'processing'
		WHERE id IN (
			SELECT id FROM sms_queue
			WHERE status = 'pending'
			ORDER BY created_at ASC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, msisdn, raw_message
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.SMSQueueItem, error) {
		var item models.SMSQueueItem
		return item, row.Scan(&item.ID, &item.MSISDN, &item.RawMessage)
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mgRepository) MarkSMSAsProcessed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sms_queue SET status = 'processed', processed_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *mgRepository) MarkSMSRejected(ctx context.Context, id uuid.UUID, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE sms_queue SET status = 'rejected', rejected_at = NOW(), reason = $2 WHERE id = $1
	`, id, reason)
	return err
}
