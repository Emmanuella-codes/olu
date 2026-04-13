package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
	"github.com/google/uuid"
)

func TestCandidateHandlerList_DBError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{allCandidatesErr: errDBDown})

	r := newTestRouter()
	r.GET("/candidates", NewCandidateHandler(services.NewCandidateService(nil, nil)).List)

	w := performRequest(r, "GET", "/candidates", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestCandidateHandlerList_Success(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{allCandidates: []models.Candidate{{Code: "A1", Name: "Ada"}}})

	r := newTestRouter()
	r.GET("/candidates", NewCandidateHandler(services.NewCandidateService(nil, nil)).List)

	w := performRequest(r, "GET", "/candidates", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}

func TestCandidateHandlerGetByID_InvalidID(t *testing.T) {
	r := newTestRouter()
	r.GET("/candidates/:id", NewCandidateHandler(services.NewCandidateService(nil, nil)).GetByID)

	w := performRequest(r, "GET", "/candidates/not-a-uuid", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCandidateHandlerGetByID_NotFound(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidate: nil})

	r := newTestRouter()
	r.GET("/candidates/:id", NewCandidateHandler(services.NewCandidateService(nil, nil)).GetByID)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body)
	}
}

func TestCandidateHandlerGetByID_DBError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidateErr: errDBDown})

	r := newTestRouter()
	r.GET("/candidates/:id", NewCandidateHandler(services.NewCandidateService(nil, nil)).GetByID)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestCandidateHandlerGetByID_Success(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidate: &models.Candidate{Code: "A1", Name: "Ada"}})

	r := newTestRouter()
	r.GET("/candidates/:id", NewCandidateHandler(services.NewCandidateService(nil, nil)).GetByID)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
