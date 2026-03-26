package vote

import (
	"context"

	"log"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository interface {
	GetAllCandidates(ctx context.Context) ([]models.Candidate, error)
	GetCandidateByID(ctx context.Context, id uuid.UUID) (*models.Candidate, error)
	GetCandidateByCode(ctx context.Context, code string) (*models.Candidate, error)
	HasVoted(ctx context.Context, voterHash string) (bool, error)
	RecordVote(ctx context.Context, vote models.Vote) error
	GetVoteTally(ctx context.Context) ([]models.TallyRow, error)
	GetTotalVoteCount(ctx context.Context) (int64, error)
	WriteAuditLog(ctx context.Context, entry models.AuditEntry)
	EnqueueSMS(ctx context.Context, msisdn, rawMessage string) error
	GetPendingSMSBatch(ctx context.Context, limit int) ([]models.SMSQueueItem, error)
	MarkSMSAsProcessed(ctx context.Context, id uuid.UUID) error
	MarkSMSRejected(ctx context.Context, id uuid.UUID, reason string) error
}

var VoteRepo VoteRepository

func InitVoteRepo(db *pgxpool.Pool, logger *log.Logger) {
	VoteRepo = newMgRepository(db, logger)
}
