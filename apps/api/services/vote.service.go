package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/repository/vote"
	"github.com/emmanuella-codes/olu/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var ErrAlreadyVoted = errors.New("voter has already voted")
var ErrInvalidCandidate = errors.New("candidate not found")

type VoteService struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
	salt string
}

func NewVoteService(pool *pgxpool.Pool, rdb *redis.Client) *VoteService {
	return &VoteService{pool: pool, rdb: rdb, salt: "olu-vote-v1"}
}

type CastVoteInput struct {
	Phone         string
	CandidateCode string
	Channel       models.VoteChannel
	IPAddress     string
	UserAgent     string
}

type CastVoteResult struct {
	ConfirmationID string
	CandidateName  string
}

func (s *VoteService) CastVote(ctx context.Context, input CastVoteInput) (*CastVoteResult, error) {
	// resolve candidate
	candidate, err := vote.VoteRepo.GetCandidateByCode(ctx, input.CandidateCode)
	if err != nil {
		return nil, fmt.Errorf("vote: lookup candidate: %w", err)
	}
	if candidate == nil {
		s.audit(ctx, input, "invalid_code")
		return nil, ErrInvalidCandidate
	}

	// hash voter identity & check for dupe
	voterHash := utils.HashVoterIdentity(input.Phone, s.salt)
	alreadyVoted, err := vote.VoteRepo.HasVoted(ctx, voterHash)
	if err != nil {
		return nil, fmt.Errorf("vote: check duplicate: %w", err)
	}
	if alreadyVoted {
		s.audit(ctx, input, "duplicate")
		return nil, ErrAlreadyVoted
	}

	// record vote
	voteModel := models.Vote{
		ID:          uuid.New(),
		CandidateID: candidate.ID,
		VoterHash:   voterHash,
		Channel:     input.Channel,
	}

	if err := vote.VoteRepo.RecordVote(ctx, voteModel); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.audit(ctx, input, "duplicate")
			return nil, ErrAlreadyVoted
		}
		s.audit(ctx, input, "error")
		return nil, fmt.Errorf("vote: record: %w", err)
	}

	s.audit(ctx, input, "success")
	log.Info().
		Str("candidate", candidate.Code).
		Str("channel", string(input.Channel)).
		Msg("vote recorded")

	return &CastVoteResult{
		ConfirmationID: voteModel.ID.String(),
		CandidateName:  candidate.Name,
	}, nil
}

func (s *VoteService) audit(ctx context.Context, input CastVoteInput, status string) {
	auditEntry := models.AuditEntry{
		VoterHash:     utils.HashVoterIdentity(input.Phone, s.salt),
		CandidateCode: input.CandidateCode,
		Channel:       string(input.Channel),
		Status:        status,
		IPAddress:     input.IPAddress,
		UserAgent:     input.UserAgent,
	}
	if err := vote.VoteRepo.WriteAuditLog(ctx, auditEntry); err != nil {
		log.Warn().Err(err).Str("status", status).Msg("audit log write failed")
	}
}
