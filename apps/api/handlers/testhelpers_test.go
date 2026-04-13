package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	voterepo "github.com/emmanuella-codes/olu/repositories/vote"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// fakeVoteRepo is a test double for the vote repository global.
type fakeVoteRepo struct {
	candidate        *models.Candidate
	candidateErr     error
	allCandidates    []models.Candidate
	allCandidatesErr error
	hasVoted         bool
	hasVotedErr      error
	recordErr        error
	tally            []models.TallyRow
	tallyErr         error
	total            int64
	totalErr         error
}

func (r *fakeVoteRepo) GetAllCandidates(_ context.Context) ([]models.Candidate, error) {
	return r.allCandidates, r.allCandidatesErr
}
func (r *fakeVoteRepo) GetCandidateByID(_ context.Context, _ uuid.UUID) (*models.Candidate, error) {
	return r.candidate, r.candidateErr
}
func (r *fakeVoteRepo) GetCandidateByCode(_ context.Context, _ string) (*models.Candidate, error) {
	return r.candidate, r.candidateErr
}
func (r *fakeVoteRepo) HasVoted(_ context.Context, _ string) (bool, error) {
	return r.hasVoted, r.hasVotedErr
}
func (r *fakeVoteRepo) RecordVote(_ context.Context, _ models.Vote) error {
	return r.recordErr
}
func (r *fakeVoteRepo) GetVoteTally(_ context.Context) ([]models.TallyRow, error) {
	return r.tally, r.tallyErr
}
func (r *fakeVoteRepo) GetTotalVoteCount(_ context.Context) (int64, error) {
	return r.total, r.totalErr
}
func (r *fakeVoteRepo) WriteAuditLog(_ context.Context, _ models.AuditEntry) error {
	return nil
}

var _ voterepo.VoteRepository = (*fakeVoteRepo)(nil)

// withFakeVoteRepo mutates the package-level VoteRepo; do not use t.Parallel in tests that call it.
func withFakeVoteRepo(t *testing.T, repo *fakeVoteRepo) {
	t.Helper()
	orig := voterepo.VoteRepo
	t.Cleanup(func() { voterepo.VoteRepo = orig })
	voterepo.VoteRepo = repo
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func performRequest(r *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder {
	var buf *bytes.Reader
	if body != nil {
		buf = bytes.NewReader(body)
	} else {
		buf = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var errDBDown = errors.New("db down")
