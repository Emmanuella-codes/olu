package adminhandler

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
)

func TestStats_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{statsErr: errDBDown})

	r := newTestRouter()
	r.GET("/stats", newHandler().Stats)

	w := performRequest(r, "GET", "/stats", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestStats_Success(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{stats: models.Stats{TotalVotes: 10, SMSVotes: 4}})

	r := newTestRouter()
	r.GET("/stats", newHandler().Stats)

	w := performRequest(r, "GET", "/stats", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
