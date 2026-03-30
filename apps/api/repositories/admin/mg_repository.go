package admin

import (
	"context"
	"log"

	"github.com/emmanuella-codes/olu/dtos"
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

func (r *mgRepository) CreateCandidate(ctx context.Context, candidate dtos.CreateCandidateDTO) (*models.Candidate, error) {
	var out models.Candidate
	id := uuid.New()
	err := r.db.QueryRow(ctx, `
		INSERT INTO candidates (id, code, name, party, bio, achievements, photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, code, name, party, bio, achievements, photo_url, is_active, created_at
	`, id, candidate.Code, candidate.Name, candidate.Party, candidate.Bio, candidate.Achievements, candidate.PhotoURL).
		Scan(&out.ID, &out.Code, &out.Name, &out.Party, &out.Bio, &out.Achievements, &out.PhotoURL, &out.IsActive, &out.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *mgRepository) UpdateCandidate(ctx context.Context, id uuid.UUID, candidate dtos.UpdateCandidateDTO) (*models.Candidate, error) {
	var out models.Candidate
	err := r.db.QueryRow(ctx, `
	UPDATE candidates
	SET code = COALESCE($2, code),
	    name = COALESCE($3, name),
	    party = COALESCE($4, party),
	    bio = COALESCE($5, bio),
	    achievements = COALESCE($6, achievements),
	    photo_url = COALESCE($7, photo_url),
	    is_active = COALESCE($8, is_active)
	WHERE id = $1
	RETURNING id, code, name, party, bio, achievements, photo_url, is_active, created_at
	`, id,
		optionalString(candidate.Code),
		optionalString(candidate.Name),
		optionalString(candidate.Party),
		optionalString(candidate.Bio),
		optionalString(candidate.Achievements),
		optionalString(candidate.PhotoURL),
		optionalBool(candidate.IsActive),
	).
		Scan(&out.ID, &out.Code, &out.Name, &out.Party, &out.Bio, &out.Achievements, &out.PhotoURL, &out.IsActive, &out.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func optionalString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}

func optionalBool(value *bool) any {
	if value == nil {
		return nil
	}
	return *value
}

func (r *mgRepository) DeactivateCandidate(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE candidates
		SET is_active = FALSE
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *mgRepository) GetAllCandidates(ctx context.Context) ([]models.Candidate, error) {
	rows, err := r.db.Query(ctx, `
	SELECT id, code, name, party, bio, achievements, photo_url, is_active, created_at
	FROM candidates
	ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Candidate
	for rows.Next() {
		var candidate models.Candidate
		if err := rows.Scan(&candidate.ID, &candidate.Code, &candidate.Name, &candidate.Party, &candidate.Bio, &candidate.Achievements, &candidate.PhotoURL, &candidate.IsActive, &candidate.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, candidate)
	}
	return out, rows.Err()
}

func (r *mgRepository) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.QueryRow(ctx, `
		SELECT id, email, password_hash, is_active, last_login
		FROM admins
		WHERE email = $1
		LIMIT 1
	`, email).Scan(&admin.ID, &admin.Email, &admin.PasswordHash, &admin.IsActive, &admin.LastLogin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &admin, nil
}

func (r *mgRepository) UpdateAdminLastLogin(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE admins
		SET last_login = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *mgRepository) GetAllStats(ctx context.Context) (models.Stats, error) {
	var stats models.Stats
	err := r.db.QueryRow(ctx, `
	SELECT 
		COUNT(*) AS total_votes,
		COUNT(*) FILTER (WHERE channel = 'web') AS web_votes,
		COUNT(*) FILTER (WHERE channel = 'sms') AS sms_votes,
		(SELECT COUNT(*) FROM sms_queue WHERE status = 'pending') AS pending_sms,
		(SELECT COUNT(*) FROM candidates) AS total_candidates
	FROM votes
	`).Scan(&stats.TotalVotes, &stats.WebVotes, &stats.SMSVotes, &stats.PendingSMS, &stats.TotalCandidates)
	if err != nil {
		return models.Stats{}, err
	}
	return stats, nil
}
