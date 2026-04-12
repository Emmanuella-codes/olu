package handlers

import (
	"net/http"
	"testing"

	"github.com/emmanuella-codes/olu/services"
)

func newOTPSvc() *services.OTPService {
	// nil rdb — VerifyCode/Send will fail at the cache layer,
	// which is what we want for the service-error path tests.
	return services.NewOTPService(nil, nil)
}

// --- maskPhone ---

func TestMaskPhone_Normal(t *testing.T) {
	got := maskPhone("2349090903080")
	want := "2349******3080"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestMaskPhone_ShortInput(t *testing.T) {
	got := maskPhone("123")
	if got != "****" {
		t.Fatalf("expected ****, got %q", got)
	}
}

// --- OTPHandler.Send ---

func TestOTPSend_MissingPhone(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/send", NewOTPHandler(newOTPSvc(), "secret").Send)

	w := performRequest(r, "POST", "/otp/send", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestOTPSend_InvalidPhone(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/send", NewOTPHandler(newOTPSvc(), "secret").Send)

	w := performRequest(r, "POST", "/otp/send", []byte(`{"phone":"not-a-phone"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

// With a nil Redis client, cache.SetOTP returns an error, causing Send → 500.
func TestOTPSend_ServiceError(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/send", NewOTPHandler(newOTPSvc(), "secret").Send)

	w := performRequest(r, "POST", "/otp/send", []byte(`{"phone":"08012345678"}`))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", w.Code, w.Body)
	}
}

// --- OTPHandler.Verify ---

func TestOTPVerify_MissingFields(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/verify", NewOTPHandler(newOTPSvc(), "secret").Verify)

	w := performRequest(r, "POST", "/otp/verify", nil)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

func TestOTPVerify_InvalidPhone(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/verify", NewOTPHandler(newOTPSvc(), "secret").Verify)

	w := performRequest(r, "POST", "/otp/verify", []byte(`{"phone":"bad","code":"123456"}`))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body)
	}
}

// With nil Redis, VerifyCode returns a retrieve error → 401.
func TestOTPVerify_ServiceError(t *testing.T) {
	r := newTestRouter()
	r.POST("/otp/verify", NewOTPHandler(newOTPSvc(), "secret").Verify)

	w := performRequest(r, "POST", "/otp/verify", []byte(`{"phone":"08012345678","code":"123456"}`))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body)
	}
}
