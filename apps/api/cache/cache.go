package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis: parse url: %w", err)
	}
	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return client, nil
}

const otpTTL = 10 * time.Minute

func OTPKey(msisdn string) string        { return fmt.Sprintf("otp:%s", msisdn) }
func OTPVerifiedKey(phone string) string { return fmt.Sprintf("otp_verified:%s", phone) }

func SetOTP(ctx context.Context, rdb *redis.Client, msisdn, code string) error {
	if rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return rdb.Set(ctx, OTPKey(msisdn), code, otpTTL).Err()
}

func GetOTP(ctx context.Context, rdb *redis.Client, msisdn string) (string, error) {
	if rdb == nil {
		return "", fmt.Errorf("redis client is nil")
	}
	val, err := rdb.Get(ctx, OTPKey(msisdn)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func DeleteOTP(ctx context.Context, rdb *redis.Client, msisdn string) error {
	if rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return rdb.Del(ctx, OTPKey(msisdn)).Err()
}

const resultsCacheKey = "results:tally"
const resultsCacheTTL = 60 * time.Second

func SetResultsCache(ctx context.Context, rdb *redis.Client, data []byte) error {
	if rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return rdb.Set(ctx, resultsCacheKey, data, resultsCacheTTL).Err()
}

func GetResultsCache(ctx context.Context, rdb *redis.Client) ([]byte, error) {
	if rdb == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	val, err := rdb.Get(ctx, resultsCacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

const candidatesCacheKey = "candidates:all"
const candidatesCacheTTL = 5 * time.Second

func SetCandidatesCache(ctx context.Context, rdb *redis.Client, data []byte) error {
	if rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	return rdb.Set(ctx, candidatesCacheKey, data, candidatesCacheTTL).Err()
}

func GetCandidatesCache(ctx context.Context, rdb *redis.Client) ([]byte, error) {
	if rdb == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	val, err := rdb.Get(ctx, candidatesCacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func RateLimitKey(prefix, identifier string) string {
	return fmt.Sprintf("rl:%s:%s", prefix, identifier)
}

func IncreaseRateLimit(ctx context.Context, rdb *redis.Client, key string, window time.Duration) (int64, error) {
	if rdb == nil {
		return 0, fmt.Errorf("redis client is nil")
	}
	pipe := rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incr.Val(), nil
}
