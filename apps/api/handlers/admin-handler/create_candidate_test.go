package adminhandler

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateCandidate_MissingFields(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	w := performRequest(r, "POST", "/candidates", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateCandidate_InvalidCode(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	body := []byte(`{"code":"INVALID!!!","name":"Ada","party":"pdp","bio":"bio","achievements":"a"}`)
	w := performRequest(r, "POST", "/candidates", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateCandidate_InvalidParty(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	body := []byte(`{"code":"A1","name":"Ada","party":"notaparty","bio":"bio","achievements":"a"}`)
	w := performRequest(r, "POST", "/candidates", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateCandidate_DuplicateCode(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{createCandidateErr: &pgconn.PgError{Code: "23505"}})

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	body := []byte(`{"code":"A1","name":"Ada","party":"pdp","bio":"bio","achievements":"a"}`)
	w := performRequest(r, "POST", "/candidates", body)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateCandidate_DBError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{createCandidateErr: errDBDown})

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	body := []byte(`{"code":"A1","name":"Ada","party":"pdp","bio":"bio","achievements":"a"}`)
	w := performRequest(r, "POST", "/candidates", body)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateCandidate_SuccessNormalizesCodeAndParty(t *testing.T) {
	repo := &fakeAdminRepo{
		createdCandidate: &models.Candidate{Code: "A1", Name: "Ada", Party: "pdp"},
	}
	withFakeAdminRepo(t, repo)

	r := newTestRouter()
	r.POST("/candidates", newHandler().CreateCandidate)

	body := []byte(`{"code":" a1 ","name":"Ada","party":" PDP ","bio":"bio","achievements":"a"}`)
	w := performRequest(r, "POST", "/candidates", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body)
	}
	if repo.createCandidateDTO.Code != "A1" {
		t.Fatalf("expected normalized code A1, got %q", repo.createCandidateDTO.Code)
	}
	if repo.createCandidateDTO.Party != "pdp" {
		t.Fatalf("expected normalized party pdp, got %q", repo.createCandidateDTO.Party)
	}
}
