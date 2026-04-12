package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newMiddlewareTestRouter(mw gin.HandlerFunc, assertions func(*gin.Context)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", mw, func(c *gin.Context) {
		if assertions != nil {
			assertions(c)
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func performMiddlewareRequest(r http.Handler, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func signOTPToken(t *testing.T, secret string, claims OTPClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign otp token: %v", err)
	}
	return token
}

func signAdminToken(t *testing.T, secret string, claims AdminClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign admin token: %v", err)
	}
	return token
}

// --- WebhookSecret ---

func TestWebhookSecretMissingHeader(t *testing.T) {
	r := newMiddlewareTestRouter(WebhookSecret("secret"), nil)

	w := performMiddlewareRequest(r, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestWebhookSecretEmptyHeader(t *testing.T) {
	r := newMiddlewareTestRouter(WebhookSecret("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"X-Webhook-Secret": ""})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for empty header value, got %d: %s", w.Code, w.Body)
	}
}

func TestWebhookSecretWrongHeader(t *testing.T) {
	r := newMiddlewareTestRouter(WebhookSecret("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"X-Webhook-Secret": "wrong"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestWebhookSecretCorrectHeader(t *testing.T) {
	r := newMiddlewareTestRouter(WebhookSecret("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"X-Webhook-Secret": "secret"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}

// An empty configured secret matches an empty header value because
// ConstantTimeCompare([]byte(""), []byte("")) == 1.
// This test documents that behaviour so it is not changed accidentally.
func TestWebhookSecretEmptyConfiguredSecret(t *testing.T) {
	r := newMiddlewareTestRouter(WebhookSecret(""), nil)

	// empty header triggers the incoming == "" guard → 401
	w := performMiddlewareRequest(r, map[string]string{"X-Webhook-Secret": ""})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for empty secret + empty header, got %d", w.Code)
	}
}

// --- RequireOTPToken ---

func TestRequireOTPTokenMissingHeader(t *testing.T) {
	r := newMiddlewareTestRouter(RequireOTPToken("secret"), nil)

	w := performMiddlewareRequest(r, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenNoBearerPrefix(t *testing.T) {
	r := newMiddlewareTestRouter(RequireOTPToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Token abc"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for non-Bearer scheme, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenInvalidToken(t *testing.T) {
	r := newMiddlewareTestRouter(RequireOTPToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer not-a-token"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenWrongSecret(t *testing.T) {
	token := signOTPToken(t, "correct-secret", OTPClaims{
		Phone:            "+2348012345678",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireOTPToken("wrong-secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong secret, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenExpired(t *testing.T) {
	token := signOTPToken(t, "secret", OTPClaims{
		Phone:            "+2348012345678",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireOTPToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenWrongAlgorithm(t *testing.T) {
	// Sign with the none algorithm — the middleware requires HMAC.
	claims := OTPClaims{
		Phone:            "+2348012345678",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign with none: %v", err)
	}
	r := newMiddlewareTestRouter(RequireOTPToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + tokenStr})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for non-HMAC algorithm, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenMalformedClaims(t *testing.T) {
	secret := "secret"
	token := signOTPToken(t, secret, OTPClaims{
		// Phone is empty — malformed claims
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireOTPToken(secret), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireOTPTokenValidTokenSetsPhone(t *testing.T) {
	secret := "secret"
	phone := "+2348012345678"
	token := signOTPToken(t, secret, OTPClaims{
		Phone:            phone,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireOTPToken(secret), func(c *gin.Context) {
		got, ok := c.Get("phone")
		if !ok || got != phone {
			t.Fatalf("expected phone %q in context, got %v", phone, got)
		}
	})

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}

// --- RequireAdminToken ---

func TestRequireAdminTokenMissingHeader(t *testing.T) {
	r := newMiddlewareTestRouter(RequireAdminToken("secret"), nil)

	w := performMiddlewareRequest(r, nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenNoBearerPrefix(t *testing.T) {
	r := newMiddlewareTestRouter(RequireAdminToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Token abc"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for non-Bearer scheme, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenInvalidToken(t *testing.T) {
	r := newMiddlewareTestRouter(RequireAdminToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer not-a-token"})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenWrongSecret(t *testing.T) {
	token := signAdminToken(t, "correct-secret", AdminClaims{
		AdminID:          uuid.New(),
		Email:            "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireAdminToken("wrong-secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for wrong secret, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenExpired(t *testing.T) {
	token := signAdminToken(t, "secret", AdminClaims{
		AdminID:          uuid.New(),
		Email:            "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireAdminToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenWrongAlgorithm(t *testing.T) {
	claims := AdminClaims{
		AdminID:          uuid.New(),
		Email:            "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign with none: %v", err)
	}
	r := newMiddlewareTestRouter(RequireAdminToken("secret"), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + tokenStr})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for non-HMAC algorithm, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenNilAdminID(t *testing.T) {
	secret := "secret"
	token := signAdminToken(t, secret, AdminClaims{
		// AdminID is zero value (uuid.Nil)
		Email:            "admin@example.com",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireAdminToken(secret), nil)

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}

func TestRequireAdminTokenValidTokenSetsContext(t *testing.T) {
	secret := "secret"
	adminID := uuid.New()
	email := "admin@example.com"
	token := signAdminToken(t, secret, AdminClaims{
		AdminID:          adminID,
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))},
	})
	r := newMiddlewareTestRouter(RequireAdminToken(secret), func(c *gin.Context) {
		gotAdminID, ok := c.Get("admin_id")
		if !ok || gotAdminID != adminID {
			t.Fatalf("expected admin_id %s in context, got %v", adminID, gotAdminID)
		}
		gotEmail, ok := c.Get("email")
		if !ok || gotEmail != email {
			t.Fatalf("expected email %s in context, got %v", email, gotEmail)
		}
	})

	w := performMiddlewareRequest(r, map[string]string{"Authorization": "Bearer " + token})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body)
	}
}

// --- NewAdminClaims ---

func TestNewAdminClaims(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	adminID := uuid.New()
	claims := NewAdminClaims(adminID, "admin@example.com", now)

	if claims.AdminID != adminID {
		t.Fatalf("expected admin id %s, got %s", adminID, claims.AdminID)
	}
	if claims.Email != "admin@example.com" {
		t.Fatalf("expected email admin@example.com, got %q", claims.Email)
	}
	if claims.Issuer != "olu-admin" {
		t.Fatalf("expected issuer olu-admin, got %q", claims.Issuer)
	}
	if claims.Subject != adminID.String() {
		t.Fatalf("expected subject %s, got %q", adminID, claims.Subject)
	}
	if !claims.IssuedAt.Time.Equal(now) {
		t.Fatalf("expected iat %s, got %s", now, claims.IssuedAt.Time)
	}
	if !claims.ExpiresAt.Time.Equal(now.Add(AdminTokenTTL)) {
		t.Fatalf("expected expiry %s, got %s", now.Add(AdminTokenTTL), claims.ExpiresAt.Time)
	}
}

// --- RateLimit ---

// With a nil Redis client, cache.IncreaseRateLimit returns an error and the
// middleware fails open — the request passes through rather than blocking.
func TestRateLimitFailsOpenOnRedisError(t *testing.T) {
	r := newMiddlewareTestRouter(RateLimit(nil, "test", 5, time.Minute), nil)

	w := performMiddlewareRequest(r, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected fail-open 200, got %d: %s", w.Code, w.Body)
	}
}
