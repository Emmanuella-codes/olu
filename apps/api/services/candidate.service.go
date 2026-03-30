package services

import (
	"context"
	"encoding/json"

	"github.com/emmanuella-codes/olu/cache"
	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/repositories/vote"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type CandidateService struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func NewCandidateService(pool *pgxpool.Pool, rdb *redis.Client) *CandidateService {
	return &CandidateService{pool: pool, rdb: rdb}
}

func (s *CandidateService) List(ctx context.Context) ([]models.Candidate, error) {
	if cached, err := cache.GetCandidatesCache(ctx, s.rdb); err == nil && cached != nil {
		var candidates []models.Candidate
		if err := json.Unmarshal(cached, &candidates); err == nil {
			return candidates, nil
		}
	}

	candidates, err := vote.VoteRepo.GetAllCandidates(ctx)
	if err != nil {
		return nil, err
	}

	// populate cache
	if data, err := json.Marshal(candidates); err == nil {
		if err := cache.SetCandidatesCache(ctx, s.rdb, data); err == nil {
			return candidates, nil
		}
	}

	return candidates, nil
}

func (s *CandidateService) GetByID(ctx context.Context, id uuid.UUID) (*models.Candidate, error) {
	return vote.VoteRepo.GetCandidateByID(ctx, id)
}
