package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	voterepo "github.com/emmanuella-codes/olu/repositories/vote"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// shared across vote and candidate service tests.
type fakeVoteRepo struct {
	candidate        *models.Candidate
	candidateErr     error
	allCandidates    []models.Candidate
	allCandidatesErr error
	hasVoted         bool
	hasVotedErr      error
	recordedVote     *models.Vote
	recordErr        error
	tally            []models.TallyRow
	tallyErr         error
	total            int64
	totalErr         error
	auditErr         error
}

func (r *fakeVoteRepo) GetAllCandidates(_ context.Context) ([]models.Candidate, error) {
	return r.allCandidates, r.allCandidatesErr
}
func (r *fakeVoteRepo) GetCandidateByID(_ context.Context, _ uuid.UUID) (*models.Candidate, error) {
	return r.candidate, r.candidateErr
}
func (r *fakeVoteRepo) GetCandidateByCode(_ context.Context, _ string) (*models.Candidate, error) {
	return r.candidate, r.candidateErr
}
func (r *fakeVoteRepo) HasVoted(_ context.Context, _ string) (bool, error) {
	return r.hasVoted, r.hasVotedErr
}
func (r *fakeVoteRepo) RecordVote(_ context.Context, vote models.Vote) error {
	r.recordedVote = &vote
	return r.recordErr
}
func (r *fakeVoteRepo) GetVoteTally(_ context.Context) ([]models.TallyRow, error) {
	return r.tally, r.tallyErr
}
func (r *fakeVoteRepo) GetTotalVoteCount(_ context.Context) (int64, error) {
	return r.total, r.totalErr
}
func (r *fakeVoteRepo) WriteAuditLog(_ context.Context, _ models.AuditEntry) error {
	return r.auditErr
}

var _ voterepo.VoteRepository = (*fakeVoteRepo)(nil)

// swaps the package-level VoteRepo for the duration of a test.
// Do not use t.Parallel with tests that call this helper.
func withFakeVoteRepo(t *testing.T, repo *fakeVoteRepo) {
	t.Helper()
	orig := voterepo.VoteRepo
	t.Cleanup(func() { voterepo.VoteRepo = orig })
	voterepo.VoteRepo = repo
}

func TestNewVoteServiceSetsDefaultSalt(t *testing.T) {
	svc := NewVoteService(nil, nil, nil)
	if svc == nil {
		t.Fatal("expected service, got nil")
	}
	if svc.salt != "olu-vote-v1" {
		t.Fatalf("expected default salt olu-vote-v1, got %q", svc.salt)
	}
}

func TestVoteServiceQueueHelpersNoopWithoutSMSService(t *testing.T) {
	svc := NewVoteService(nil, nil, nil)

	svc.queueConfirmation(context.Background(), "09090903080", "Ada", "confirmation-id")
	svc.queueRejection(context.Background(), "09090903080", "rejected")
}

func TestCastVote_InvalidCandidate(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidate: nil})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "Z9", Channel: "web",
	})
	if !errors.Is(err, ErrInvalidCandidate) {
		t.Fatalf("expected ErrInvalidCandidate, got %v", err)
	}
}

func TestCastVote_CandidateLookupError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidateErr: errors.New("lookup failed")})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if err == nil || !strings.Contains(err.Error(), "vote: lookup candidate") {
		t.Fatalf("expected lookup candidate error, got %v", err)
	}
}

func TestCastVote_AlreadyVoted(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		hasVoted:  true,
	})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if !errors.Is(err, ErrAlreadyVoted) {
		t.Fatalf("expected ErrAlreadyVoted, got %v", err)
	}
}

func TestCastVote_HasVotedError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate:   &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		hasVotedErr: errors.New("has voted failed"),
	})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if err == nil || !strings.Contains(err.Error(), "vote: check duplicate") {
		t.Fatalf("expected check duplicate error, got %v", err)
	}
}

func TestCastVote_DuplicateConstraint(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		recordErr: &pgconn.PgError{Code: "23505"},
	})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if !errors.Is(err, ErrAlreadyVoted) {
		t.Fatalf("expected ErrAlreadyVoted on pg 23505, got %v", err)
	}
}

func TestCastVote_RecordError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		recordErr: errors.New("db unavailable"),
	})

	svc := NewVoteService(nil, nil, nil)
	_, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if err == nil || errors.Is(err, ErrAlreadyVoted) || errors.Is(err, ErrInvalidCandidate) {
		t.Fatalf("expected a generic record error, got %v", err)
	}
}

func TestCastVote_Success(t *testing.T) {
	candidateID := uuid.New()
	repo := &fakeVoteRepo{
		candidate: &models.Candidate{ID: candidateID, Code: "A1", Name: "Ada Obi"},
	}
	withFakeVoteRepo(t, repo)

	svc := NewVoteService(nil, nil, nil)
	result, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.CandidateName != "Ada Obi" {
		t.Fatalf("expected candidate name Ada Obi, got %q", result.CandidateName)
	}
	if result.ConfirmationID == "" {
		t.Fatal("expected non-empty confirmation ID")
	}
	if repo.recordedVote == nil {
		t.Fatal("expected recorded vote")
	}
	if repo.recordedVote.CandidateID != candidateID {
		t.Fatalf("expected candidate id %s, got %s", candidateID, repo.recordedVote.CandidateID)
	}
	if repo.recordedVote.Channel != models.WebVoteChannel {
		t.Fatalf("expected web channel, got %q", repo.recordedVote.Channel)
	}
	if repo.recordedVote.VoterHash == "" {
		t.Fatal("expected non-empty voter hash")
	}
}

func TestCastVote_AuditFailureDoesNotFailVote(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada Obi"},
		auditErr:  errors.New("audit failed"),
	})

	svc := NewVoteService(nil, nil, nil)
	result, err := svc.CastVote(context.Background(), CastVoteInput{
		Phone: "09090903080", CandidateCode: "A1", Channel: "web",
	})
	if err != nil {
		t.Fatalf("expected audit failure to be non-fatal, got %v", err)
	}
	if result == nil || result.CandidateName != "Ada Obi" {
		t.Fatalf("unexpected result: %+v", result)
	}
}
