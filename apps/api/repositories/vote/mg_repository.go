package vote

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

func (r *mgRepository) GetAllCandidates(ctx context.Context) ([]models.Candidate, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, code, name, party, bio, achievements, photo_url, is_active, created_at
		FROM candidates
		WHERE is_active = TRUE
		ORDER BY name DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Candidate
	for rows.Next() {
		var c models.Candidate
		if err := rows.Scan(&c.ID, &c.Code, &c.Name, &c.Party, &c.Bio, &c.Achievements, &c.PhotoURL, &c.IsActive, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *mgRepository) GetCandidateByID(ctx context.Context, id uuid.UUID) (*models.Candidate, error) {
	var c models.Candidate
	err := r.db.QueryRow(ctx, `
		SELECT id, code, name, party, bio, achievements, photo_url, is_active, created_at
		FROM candidates
		WHERE id = $1
		AND is_active = TRUE
		LIMIT 1
	`, id).Scan(&c.ID, &c.Code, &c.Name, &c.Party, &c.Bio, &c.Achievements, &c.PhotoURL, &c.IsActive, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &c, nil
}

func (r *mgRepository) GetCandidateByCode(ctx context.Context, code string) (*models.Candidate, error) {
	var c models.Candidate
	err := r.db.QueryRow(ctx, `
		SELECT id, code, name, party, bio, achievements, photo_url, is_active, created_at
		FROM candidates
		WHERE code = $1
		AND is_active = TRUE
		LIMIT 1
	`, code).Scan(&c.ID, &c.Code, &c.Name, &c.Party, &c.Bio, &c.Achievements, &c.PhotoURL, &c.IsActive, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &c, nil
}

func (r *mgRepository) HasVoted(ctx context.Context, voterHash string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
	SELECT EXISTS(SELECT 1 FROM votes WHERE voter_hash = $1)
	`, voterHash).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *mgRepository) RecordVote(ctx context.Context, vote models.Vote) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO votes (id, candidate_id, voter_hash, channel, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, vote.ID, vote.CandidateID, vote.VoterHash, vote.Channel); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *mgRepository) GetVoteTally(ctx context.Context) ([]models.TallyRow, error) {
	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.code, c.name, c.party, COUNT(v.candidate_id) AS vote_count
		FROM candidates c
		LEFT JOIN votes v ON v.candidate_id = c.id
		WHERE c.is_active = TRUE
		GROUP BY c.id, c.code, c.name, c.party
		ORDER BY vote_count DESC, c.name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tally []models.TallyRow
	for rows.Next() {
		var row models.TallyRow
		if err := rows.Scan(&row.CandidateID, &row.Code, &row.Name, &row.Party, &row.VoteCount); err != nil {
			return nil, err
		}
		tally = append(tally, row)
	}
	return tally, rows.Err()
}

func (r *mgRepository) GetTotalVoteCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM votes
	`).Scan(&count)
	return count, err
}

func (r *mgRepository) WriteAuditLog(ctx context.Context, entry models.AuditEntry) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_logs (voter_hash, candidate_code, channel, status, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, entry.VoterHash, entry.CandidateCode, entry.Channel, entry.Status, entry.IPAddress, entry.UserAgent)
	return err
}

func (r *mgRepository) EnqueueSMS(ctx context.Context, msisdn, rawMessage string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sms_queue (msisdn, message, status, created_at)
		VALUES ($1, $2, 'pending', NOW())
	`, msisdn, rawMessage)
	return err
}

func (r *mgRepository) GetPendingSMSBatch(ctx context.Context, limit int) ([]models.SMSQueueItem, error) {
	rows, err := r.db.Query(ctx, `
	SELECT id, msisdn, raw_message
	FROM sms_queue
	WHERE status = 'pending'
	ORDER BY created_at ASC
	LIMIT $1
	FOR UPDATE SKIP LOCKED
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.SMSQueueItem
	for rows.Next() {
		var item models.SMSQueueItem
		if err := rows.Scan(&item.ID, &item.MSISDN, &item.RawMessage); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
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
