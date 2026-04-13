package services

import (
	"context"
	"errors"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
)

func TestNewCandidateServiceReturnsService(t *testing.T) {
	svc := NewCandidateService(nil, nil)
	if svc == nil {
		t.Fatal("expected service, got nil")
	}
}

// with nil rdb the cache check errors and is skipped; the service falls
// through to the repository. This exercises the cache-miss → DB-success path.
func TestCandidateServiceList_CacheMissDBSuccess(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		allCandidates: []models.Candidate{
			{Code: "A1", Name: "Ada"},
		},
	})

	svc := NewCandidateService(nil, nil)
	candidates, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(candidates) != 1 || candidates[0].Code != "A1" {
		t.Fatalf("unexpected candidates: %+v", candidates)
	}
}

func TestCandidateServiceList_DBError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		allCandidatesErr: errors.New("db down"),
	})

	svc := NewCandidateService(nil, nil)
	_, err := svc.List(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCandidateServiceList_EmptyDB(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{allCandidates: []models.Candidate{}})

	svc := NewCandidateService(nil, nil)
	candidates, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(candidates) != 0 {
		t.Fatalf("expected empty slice, got %+v", candidates)
	}
}

func TestCandidateServiceGetByID_Success(t *testing.T) {
	id := uuid.New()
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: id, Code: "A1", Name: "Ada"},
	})

	svc := NewCandidateService(nil, nil)
	candidate, err := svc.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if candidate == nil || candidate.ID != id || candidate.Code != "A1" {
		t.Fatalf("unexpected candidate: %+v", candidate)
	}
}

func TestCandidateServiceGetByID_RepoError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidateErr: errors.New("lookup failed")})

	svc := NewCandidateService(nil, nil)
	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
