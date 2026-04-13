package adminhandler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestUpdateCandidate_InvalidUUID(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	w := performRequest(r, "PUT", "/candidates/not-a-uuid", []byte(`{}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_InvalidCode(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	body := []byte(`{"code":"INVALID!!!"}`)
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_InvalidParty(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	body := []byte(`{"party":"notaparty"}`)
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_NotFound(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{updatedCandidate: nil})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), []byte(`{"name":"Ada"}`))
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_DuplicateCode(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{updateCandidateErr: &pgconn.PgError{Code: "23505"}})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), []byte(`{"name":"Ada"}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{updateCandidateErr: errDBDown})

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), []byte(`{"name":"Ada"}`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestUpdateCandidate_SuccessNormalizesCodeAndParty(t *testing.T) {
	repo := &fakeAdminRepo{
		updatedCandidate: &models.Candidate{Code: "A1", Name: "Ada Updated"},
	}
	withFakeAdminRepo(t, repo)

	r := newTestRouter()
	r.PUT("/candidates/:id", newHandler().UpdateCandidate)

	id := uuid.New()
	w := performRequest(r, "PUT", fmt.Sprintf("/candidates/%s", id), []byte(`{"code":" a1 ","party":" PDP ","name":"Ada Updated"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
	if repo.updateCandidateDTO.Code == nil || *repo.updateCandidateDTO.Code != "A1" {
		t.Fatalf("expected normalized code A1, got %+v", repo.updateCandidateDTO.Code)
	}
	if repo.updateCandidateDTO.Party == nil || *repo.updateCandidateDTO.Party != "pdp" {
		t.Fatalf("expected normalized party pdp, got %+v", repo.updateCandidateDTO.Party)
	}
}
