package adminhandler

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
)

func TestAllCandidates_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{allCandidatesErr: errDBDown})

	r := newTestRouter()
	r.GET("/candidates", newHandler().AllCandidates)

	w := performRequest(r, "GET", "/candidates", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestAllCandidates_Success(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{allCandidates: []models.Candidate{{Code: "A1", Name: "Ada"}}})

	r := newTestRouter()
	r.GET("/candidates", newHandler().AllCandidates)

	w := performRequest(r, "GET", "/candidates", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
