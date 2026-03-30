package services

import "github.com/redis/go-redis/v9"

type OTPService struct {
	rdb *redis.Client
	// sms provider
}

// func NewOTPService(rdb *redis.Client) *OTPService {}
