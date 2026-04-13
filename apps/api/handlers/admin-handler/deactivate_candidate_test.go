package adminhandler

import (
	"fmt"
	"net/http"
	"testing"

	adminrepo "github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/google/uuid"
)

func TestDeactivateCandidate_InvalidUUID(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.DELETE("/candidates/:id", newHandler().DeactivateCandidate)

	w := performRequest(r, "DELETE", "/candidates/not-a-uuid", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestDeactivateCandidate_NotFound(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{deactivateErr: adminrepo.ErrNotFound})

	r := newTestRouter()
	r.DELETE("/candidates/:id", newHandler().DeactivateCandidate)

	id := uuid.New()
	w := performRequest(r, "DELETE", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body)
	}
}

func TestDeactivateCandidate_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{deactivateErr: errDBDown})

	r := newTestRouter()
	r.DELETE("/candidates/:id", newHandler().DeactivateCandidate)

	id := uuid.New()
	w := performRequest(r, "DELETE", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestDeactivateCandidate_Success(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{deactivateErr: nil})

	r := newTestRouter()
	r.DELETE("/candidates/:id", newHandler().DeactivateCandidate)

	id := uuid.New()
	w := performRequest(r, "DELETE", fmt.Sprintf("/candidates/%s", id), nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}
