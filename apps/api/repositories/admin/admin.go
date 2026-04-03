package admin

import (
	"context"
	"errors"
	"log"

	"github.com/emmanuella-codes/olu/dtos"
	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type AdminRepository interface {
	CreateAdmin(ctx context.Context, email, passwordHash string) (*models.Admin, error)
	CreateCandidate(ctx context.Context, candidate dtos.CreateCandidateDTO) (*models.Candidate, error)
	UpdateCandidate(ctx context.Context, id uuid.UUID, candidate dtos.UpdateCandidateDTO) (*models.Candidate, error)
	DeactivateCandidate(ctx context.Context, id uuid.UUID) error
	GetAllCandidates(ctx context.Context) ([]models.Candidate, error)
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
	UpdateAdminLastLogin(ctx context.Context, id uuid.UUID) error
	GetAllStats(ctx context.Context) (models.Stats, error)
}

var AdminRepo AdminRepository

func InitAdminRepo(db *pgxpool.Pool, logger *log.Logger) {
	AdminRepo = newMgRepository(db, logger)
}
