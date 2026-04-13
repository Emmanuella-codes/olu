package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/emmanuella-codes/olu/models"
)

func TestLeadersFromTallyReturnsNilForEmptyTally(t *testing.T) {
	leaders := leadersFromTally(nil)
	if leaders != nil {
		t.Fatalf("expected nil leaders, got %+v", leaders)
	}
}

func TestLeadersFromTallyReturnsSingleCandidate(t *testing.T) {
	tally := []models.TallyRow{
		{Name: "Ada", VoteCount: 7},
	}
	leaders := leadersFromTally(tally)
	if len(leaders) != 1 || leaders[0].Name != "Ada" {
		t.Fatalf("expected Ada as only leader, got %+v", leaders)
	}
}

func TestLeadersFromTallyReturnsSingleLeader(t *testing.T) {
	tally := []models.TallyRow{
		{Name: "Ada", VoteCount: 5},
		{Name: "Ben", VoteCount: 3},
	}

	leaders := leadersFromTally(tally)
	if len(leaders) != 1 || leaders[0].Name != "Ada" {
		t.Fatalf("expected Ada as only leader, got %+v", leaders)
	}
}

func TestLeadersFromTallyReturnsTieLeaders(t *testing.T) {
	tally := []models.TallyRow{
		{Name: "Ada", VoteCount: 5},
		{Name: "Ben", VoteCount: 5},
		{Name: "Chi", VoteCount: 3},
	}

	leaders := leadersFromTally(tally)
	if len(leaders) != 2 {
		t.Fatalf("expected 2 leaders, got %+v", leaders)
	}
	if leaders[0].Name != "Ada" || leaders[1].Name != "Ben" {
		t.Fatalf("unexpected leaders: %+v", leaders)
	}
}

func TestLeadersFromTallyAllTied(t *testing.T) {
	tally := []models.TallyRow{
		{Name: "Ada", VoteCount: 3},
		{Name: "Ben", VoteCount: 3},
		{Name: "Chi", VoteCount: 3},
	}
	leaders := leadersFromTally(tally)
	if len(leaders) != 3 {
		t.Fatalf("expected all 3 candidates as leaders, got %+v", leaders)
	}
}

func TestResultsServiceGetResults_SuccessWithTie(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		tally: []models.TallyRow{
			{Name: "Ada", VoteCount: 5},
			{Name: "Ben", VoteCount: 5},
			{Name: "Chi", VoteCount: 3},
		},
		total: 13,
	})

	svc := NewResultsService(nil, nil)
	results, err := svc.GetResults(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if results.TotalVotes != 13 {
		t.Fatalf("expected total votes 13, got %d", results.TotalVotes)
	}
	if !results.IsTie {
		t.Fatal("expected tie")
	}
	if len(results.Leaders) != 2 {
		t.Fatalf("expected 2 leaders, got %+v", results.Leaders)
	}
}

func TestResultsServiceGetResults_TallyError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{tallyErr: errors.New("tally failed")})

	svc := NewResultsService(nil, nil)
	_, err := svc.GetResults(context.Background())
	if err == nil || !strings.Contains(err.Error(), "results: get tally") {
		t.Fatalf("expected get tally error, got %v", err)
	}
}

func TestResultsServiceGetResults_TotalCountError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		tally:    []models.TallyRow{{Name: "Ada", VoteCount: 5}},
		totalErr: errors.New("count failed"),
	})

	svc := NewResultsService(nil, nil)
	_, err := svc.GetResults(context.Background())
	if err == nil || !strings.Contains(err.Error(), "results: get total vote count") {
		t.Fatalf("expected total vote count error, got %v", err)
	}
}
