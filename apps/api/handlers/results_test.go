package handlers

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
)

func TestResultsHandler_DBError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{tallyErr: errDBDown})

	r := newTestRouter()
	r.GET("/results", NewResultsHandler(services.NewResultsService(nil, nil)).GetResults)

	w := performRequest(r, "GET", "/results", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestResultsHandler_Success(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		tally: []models.TallyRow{{Name: "Ada", VoteCount: 3}},
		total: 3,
	})

	r := newTestRouter()
	r.GET("/results", NewResultsHandler(services.NewResultsService(nil, nil)).GetResults)

	w := performRequest(r, "GET", "/results", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
