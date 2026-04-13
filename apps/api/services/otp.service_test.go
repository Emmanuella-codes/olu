package services

import (
	"context"
	"strings"
	"testing"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateOTPReturnsRequestedDigits(t *testing.T) {
	code, err := generateOTP(6)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected 6 digits, got %q", code)
	}
	for _, ch := range code {
		if !unicode.IsDigit(ch) {
			t.Fatalf("expected only digits, got %q", code)
		}
	}
}

// run many iterations to catch cases where a leading-zero OTP would be
// shorter than expected if the format string were wrong.
func TestGenerateOTPAlwaysProducesFullWidth(t *testing.T) {
	for i := range 200 {
		code, err := generateOTP(6)
		if err != nil {
			t.Fatalf("iteration %d: unexpected error: %v", i, err)
		}
		if len(code) != 6 {
			t.Fatalf("iteration %d: expected 6 chars, got %q (len=%d)", i, code, len(code))
		}
	}
}

func TestGenerateOTPPreservesLeadingZeroWidth(t *testing.T) {
	code, err := generateOTP(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(strings.TrimSpace(code)) != 1 {
		t.Fatalf("expected one digit, got %q", code)
	}
}

func TestIssueOTPTokenIncludesPhoneClaim(t *testing.T) {
	secret := "test-secret"
	tokenString, err := issueOTPToken("09090903080", secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if !token.Valid {
		t.Fatal("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected map claims, got %T", token.Claims)
	}
	if claims["phone"] != "09090903080" {
		t.Fatalf("expected phone claim 09090903080, got %v", claims["phone"])
	}
	if _, ok := claims["exp"]; !ok {
		t.Fatal("expected exp claim")
	}
}

// calls cache.GetOTP which requires Redis. With a nil client,
// the cache layer returns an error, exercising the retrieve-error path.
func TestVerifyCode_NilRedisReturnsRetrieveError(t *testing.T) {
	svc := &OTPService{rdb: nil}
	_, err := svc.VerifyCode(context.Background(), "09090903080", "123456", "secret")
	if err == nil || !strings.Contains(err.Error(), "otp: retrieve") {
		t.Fatalf("expected otp: retrieve error, got %v", err)
	}
}
