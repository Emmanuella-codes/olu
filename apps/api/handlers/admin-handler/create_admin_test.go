package adminhandler

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateAdmin_MissingFields(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_ShortPassword(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{admin: nil})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":"a@b.com","password":"short"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_EmailAlreadyExists(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin: &models.Admin{ID: uuid.New(), Email: "a@b.com"},
	})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":"a@b.com","password":"validpassword123!"}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_LookupError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{adminErr: errDBDown})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":"a@b.com","password":"validpassword123!"}`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_CreateError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin:          nil,
		createAdminErr: errDBDown,
	})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":"a@b.com","password":"validpassword123!"}`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_DuplicateConstraint(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin:          nil,
		createAdminErr: &pgconn.PgError{Code: "23505"},
	})

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":"a@b.com","password":"validpassword123!"}`))
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body)
	}
}

func TestCreateAdmin_Success(t *testing.T) {
	newID := uuid.New()
	repo := &fakeAdminRepo{
		admin:        nil,
		createdAdmin: &models.Admin{ID: newID, Email: "a@b.com"},
	}
	withFakeAdminRepo(t, repo)

	r := newTestRouter()
	r.POST("/admin", newHandler().CreateAdmin)

	w := performRequest(r, "POST", "/admin", []byte(`{"email":" ADMIN@Example.COM ","password":"validpassword123!"}`))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body)
	}
	if repo.createdEmail != "admin@example.com" {
		t.Fatalf("expected normalized email admin@example.com, got %q", repo.createdEmail)
	}
	if repo.createdHash == "validpassword123!" {
		t.Fatal("expected hashed password, got plain password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(repo.createdHash), []byte("validpassword123!")); err != nil {
		t.Fatalf("expected stored password hash to match password, got %v", err)
	}
}
