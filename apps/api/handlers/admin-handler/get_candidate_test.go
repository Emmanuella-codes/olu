package adminhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
)

func TestGetCandidate_InvalidUUID(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.GET("/candidates/:id", newHandler().GetCandidate)

	w := performRequest(r, "GET", "/candidates/not-a-uuid", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestGetCandidate_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{candidateErr: errDBDown})

	r := newTestRouter()
	r.GET("/candidates/:id", newHandler().GetCandidate)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestGetCandidate_NotFound(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{candidate: nil})

	r := newTestRouter()
	r.GET("/candidates/:id", newHandler().GetCandidate)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body)
	}
}

func TestGetCandidate_Success(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{candidate: &models.Candidate{Code: "A1", Name: "Ada"}})

	r := newTestRouter()
	r.GET("/candidates/:id", newHandler().GetCandidate)

	id := uuid.New()
	w := performRequest(r, "GET", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
