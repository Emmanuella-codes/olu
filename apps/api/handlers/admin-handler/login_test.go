package adminhandler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/middleware"
	"github.com/emmanuella-codes/olu/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --- normalizeAdminEmail ---

func TestNormalizeAdminEmail_TrimsAndLowercases(t *testing.T) {
	got := normalizeAdminEmail("  ADMIN@Test.COM  ")
	if got != "admin@test.com" {
		t.Fatalf("expected admin@test.com, got %q", got)
	}
}

func TestNormalizeAdminEmail_EmptyAfterTrim(t *testing.T) {
	got := normalizeAdminEmail("   ")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

// --- AdminHandler.Login ---

func TestAdminLogin_MissingFields(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestAdminLogin_LookupError(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{adminErr: errDBDown})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", []byte(`{"email":"x@x.com","password":"secret"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestAdminLogin_AdminNotFound(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{admin: nil})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", []byte(`{"email":"x@x.com","password":"secret"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestAdminLogin_InactiveAdmin(t *testing.T) {
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin: &models.Admin{ID: uuid.New(), Email: "a@b.com", IsActive: false},
	})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", []byte(`{"email":"a@b.com","password":"secret"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestAdminLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.MinCost)
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin: &models.Admin{ID: uuid.New(), Email: "a@b.com", PasswordHash: string(hash), IsActive: true},
	})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", []byte(`{"email":"a@b.com","password":"wrongpassword"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestAdminLogin_Success(t *testing.T) {
	const password = "correctpassword"
	adminID := uuid.New()
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	withFakeAdminRepo(t, &fakeAdminRepo{
		admin: &models.Admin{ID: adminID, Email: "a@b.com", PasswordHash: string(hash), IsActive: true},
	})

	r := newTestRouter()
	r.POST("/login", newHandler().Login)

	w := performRequest(r, "POST", "/login", []byte(`{"email":"a@b.com","password":"correctpassword"}`))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}

	var payload struct {
		Token     string `json:"token"`
		ExpiresIn int    `json:"expires_in"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Token == "" {
		t.Fatal("expected token in response")
	}
	if payload.ExpiresIn != int(middleware.AdminTokenTTL.Seconds()) {
		t.Fatalf("expected expires_in %d, got %d", int(middleware.AdminTokenTTL.Seconds()), payload.ExpiresIn)
	}

	token, err := jwt.ParseWithClaims(payload.Token, &middleware.AdminClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("test-jwt-secret"), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("expected valid admin token, token=%v err=%v", token, err)
	}
	claims := token.Claims.(*middleware.AdminClaims)
	if claims.AdminID != adminID {
		t.Fatalf("expected admin id %s, got %s", adminID, claims.AdminID)
	}
	if claims.Email != "a@b.com" {
		t.Fatalf("expected email a@b.com, got %q", claims.Email)
	}
}
