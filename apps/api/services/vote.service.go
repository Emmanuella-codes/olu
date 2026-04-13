package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/repositories/vote"
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
	pool   *pgxpool.Pool
	rdb    *redis.Client
	smsSvc *SMSService
	salt   string
}

func NewVoteService(pool *pgxpool.Pool, rdb *redis.Client, smsSvc *SMSService) *VoteService {
	return &VoteService{pool: pool, rdb: rdb, smsSvc: smsSvc, salt: "olu-vote-v1"}
}

type CastVoteInput struct {
	Phone         string
	CandidateCode string
	Channel       models.VoteChannel
	IPAddress     string
	UserAgent     string
}

type CastVoteResult struct {
	ConfirmationID string `json:"confirmation_id"`
	CandidateName  string `json:"candidate_name"`
}

func (s *VoteService) CastVote(ctx context.Context, input CastVoteInput) (*CastVoteResult, error) {
	// resolve candidate
	candidate, err := vote.VoteRepo.GetCandidateByCode(ctx, input.CandidateCode)
	if err != nil {
		return nil, fmt.Errorf("vote: lookup candidate: %w", err)
	}
	if candidate == nil {
		s.audit(ctx, input, "invalid_code")
		s.queueRejection(ctx, input.Phone, "We could not record your vote because the candidate code is invalid.")
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
		s.queueRejection(ctx, input.Phone, "We could not record your vote because this number has already voted.")
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
			s.queueRejection(ctx, input.Phone, "We could not record your vote because this number has already voted.")
			return nil, ErrAlreadyVoted
		}
		s.audit(ctx, input, "error")
		s.queueRejection(ctx, input.Phone, "We could not process your vote at this time. Please try again later.")
		return nil, fmt.Errorf("vote: record: %w", err)
	}

	s.audit(ctx, input, "success")
	s.queueConfirmation(ctx, input.Phone, candidate.Name, voteModel.ID.String())
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

func (s *VoteService) queueConfirmation(ctx context.Context, phone, candidateName, confirmationID string) {
	if s.smsSvc == nil {
		return
	}
	if err := s.smsSvc.QueueVoteConfirmation(ctx, phone, candidateName, confirmationID); err != nil {
		log.Warn().Err(err).Str("phone", phone).Msg("queue vote confirmation failed")
	}
}

func (s *VoteService) queueRejection(ctx context.Context, phone, reason string) {
	if s.smsSvc == nil {
		return
	}
	if err := s.smsSvc.QueueVoteRejection(ctx, phone, reason); err != nil {
		log.Warn().Err(err).Str("phone", phone).Msg("queue vote rejection failed")
	}
}
