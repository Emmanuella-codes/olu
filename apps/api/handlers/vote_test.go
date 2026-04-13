package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func newVoteSvc() *services.VoteService {
	return services.NewVoteService(nil, nil, nil)
}

func TestVoteCast_MissingBody(t *testing.T) {
	r := newTestRouter()
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_InvalidCandidateCode(t *testing.T) {
	r := newTestRouter()
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"INVALID!!!"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_MissingPhoneInContext(t *testing.T) {
	r := newTestRouter()
	// phone is not set in context — no middleware
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_InvalidCandidateFromService(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{candidate: nil})

	r := newTestRouter()
	r.Use(func(c *gin.Context) { c.Set("phone", "2349090903080"); c.Next() })
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_AlreadyVoted(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		hasVoted:  true,
	})

	r := newTestRouter()
	r.Use(func(c *gin.Context) { c.Set("phone", "2349090903080"); c.Next() })
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_InvalidPhoneInContext(t *testing.T) {
	r := newTestRouter()
	r.Use(func(c *gin.Context) { c.Set("phone", 12345); c.Next() })
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_ServiceError(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada"},
		recordErr: errors.New("db down"),
	})

	r := newTestRouter()
	r.Use(func(c *gin.Context) { c.Set("phone", "2349090903080"); c.Next() })
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestVoteCast_Success(t *testing.T) {
	withFakeVoteRepo(t, &fakeVoteRepo{
		candidate: &models.Candidate{ID: uuid.New(), Code: "A1", Name: "Ada Obi"},
	})

	r := newTestRouter()
	r.Use(func(c *gin.Context) { c.Set("phone", "2349090903080"); c.Next() })
	r.POST("/vote", NewVoteHandler(newVoteSvc()).Cast)

	w := performRequest(r, "POST", "/vote", []byte(`{"candidate_code":"A1"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}

	var payload struct {
		Data struct {
			ConfirmationID string `json:"confirmation_id"`
			CandidateName  string `json:"candidate_name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Data.CandidateName != "Ada Obi" {
		t.Fatalf("expected candidate name Ada Obi, got %q", payload.Data.CandidateName)
	}
	if payload.Data.ConfirmationID == "" {
		t.Fatal("expected confirmation id")
	}
}
