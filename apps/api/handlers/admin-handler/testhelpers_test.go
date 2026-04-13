package adminhandler

import (
	"bytes"
	"context"
	"errors"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emmanuella-codes/olu/dtos"
	"github.com/emmanuella-codes/olu/models"
	adminrepo "github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// fakeAdminRepo is a test double for the admin repository global.
type fakeAdminRepo struct {
	admin              *models.Admin
	adminErr           error
	createdAdmin       *models.Admin
	createAdminErr     error
	createdEmail       string
	createdHash        string
	candidate          *models.Candidate
	candidateErr       error
	allCandidates      []models.Candidate
	allCandidatesErr   error
	createdCandidate   *models.Candidate
	createCandidateErr error
	createCandidateDTO dtos.CreateCandidateDTO
	updatedCandidate   *models.Candidate
	updateCandidateErr error
	updateCandidateDTO dtos.UpdateCandidateDTO
	deactivateErr      error
	stats              models.Stats
	statsErr           error
}

func (r *fakeAdminRepo) GetAdminByEmail(_ context.Context, _ string) (*models.Admin, error) {
	return r.admin, r.adminErr
}
func (r *fakeAdminRepo) CreateAdmin(_ context.Context, email, passwordHash string) (*models.Admin, error) {
	r.createdEmail = email
	r.createdHash = passwordHash
	return r.createdAdmin, r.createAdminErr
}
func (r *fakeAdminRepo) GetAllCandidates(_ context.Context) ([]models.Candidate, error) {
	return r.allCandidates, r.allCandidatesErr
}
func (r *fakeAdminRepo) GetCandidateByID(_ context.Context, _ uuid.UUID) (*models.Candidate, error) {
	return r.candidate, r.candidateErr
}
func (r *fakeAdminRepo) CreateCandidate(_ context.Context, candidate dtos.CreateCandidateDTO) (*models.Candidate, error) {
	r.createCandidateDTO = candidate
	return r.createdCandidate, r.createCandidateErr
}
func (r *fakeAdminRepo) UpdateCandidate(_ context.Context, _ uuid.UUID, candidate dtos.UpdateCandidateDTO) (*models.Candidate, error) {
	r.updateCandidateDTO = candidate
	return r.updatedCandidate, r.updateCandidateErr
}
func (r *fakeAdminRepo) DeactivateCandidate(_ context.Context, _ uuid.UUID) error {
	return r.deactivateErr
}
func (r *fakeAdminRepo) UpdateAdminLastLogin(_ context.Context, _ uuid.UUID) error {
	return nil
}
func (r *fakeAdminRepo) GetAllStats(_ context.Context) (models.Stats, error) {
	return r.stats, r.statsErr
}

var _ adminrepo.AdminRepository = (*fakeAdminRepo)(nil)

// withFakeAdminRepo mutates the package-level AdminRepo; do not use t.Parallel in tests that call it.
func withFakeAdminRepo(t *testing.T, repo *fakeAdminRepo) {
	t.Helper()
	orig := adminrepo.AdminRepo
	t.Cleanup(func() { adminrepo.AdminRepo = orig })
	adminrepo.AdminRepo = repo
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

func newHandler() *AdminHandler {
	return NewAdminHandler("test-jwt-secret")
}

var errDBDown = errors.New("db down")

// TestMain primes adminrepo.AdminRepo with a safe no-op implementation so that
// any goroutine started by a handler (e.g. the async UpdateAdminLastLogin in
// Login) never calls a nil interface after test cleanup restores the global.
func TestMain(m *testing.M) {
	adminrepo.AdminRepo = &fakeAdminRepo{}
	os.Exit(m.Run())
}
