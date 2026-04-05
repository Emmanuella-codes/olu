package services

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/emmanuella-codes/olu/cache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type OTPService struct {
	rdb    *redis.Client
	smsSvc *SMSService
}

func NewOTPService(rdb *redis.Client, smsSvc *SMSService) *OTPService {
	return &OTPService{rdb: rdb, smsSvc: smsSvc}
}

func (s *OTPService) Send(ctx context.Context, phone string) error {
	return s.RequestCode(ctx, phone)
}

func (s *OTPService) RequestCode(ctx context.Context, phone string) error {
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("otp: generate code: %w", err)
	}

	if err := cache.SetOTP(ctx, s.rdb, phone, code); err != nil {
		return fmt.Errorf("otp: store code: %w", err)
	}

	if err := s.smsSvc.SendOTP(ctx, phone, code); err != nil {
		_ = cache.DeleteOTP(ctx, s.rdb, phone)
		return fmt.Errorf("otp: send code: %w", err)
	}

	return nil
}

func (s *OTPService) VerifyCode(ctx context.Context, phone, code, jwtSecret string) (string, error) {
	stored, err := cache.GetOTP(ctx, s.rdb, phone)
	if err != nil {
		return "", fmt.Errorf("otp: retrieve: %w", err)
	}
	if stored == "" {
		return "", fmt.Errorf("otp: code expired or not found")
	}
	if subtle.ConstantTimeCompare([]byte(stored), []byte(strings.TrimSpace(code))) != 1 {
		return "", fmt.Errorf("otp: invalid code")
	}

	if err := cache.DeleteOTP(ctx, s.rdb, phone); err != nil {
		return "", fmt.Errorf("otp: delete code: %w", err)
	}

	token, err := issueOTPToken(phone, jwtSecret)
	if err != nil {
		return "", fmt.Errorf("otp: issue token: %w", err)
	}
	return token, nil
}

// short-lived JWT (15 min) proving OTP was verified.
func issueOTPToken(phone, secret string) (string, error) {
	claims := jwt.MapClaims{
		"phone": phone,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateOTP(digits int) (string, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", digits, n), nil
}
