package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/emmanuella-codes/olu/cache"
	"github.com/emmanuella-codes/olu/models"
	"github.com/emmanuella-codes/olu/repositories/vote"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ResultsService struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func NewResultsService(pool *pgxpool.Pool, rdb *redis.Client) *ResultsService {
	return &ResultsService{pool: pool, rdb: rdb}
}

func (s *ResultsService) GetResults(ctx context.Context) (*models.Results, error) {
	if cached, err := cache.GetResultsCache(ctx, s.rdb); err == nil && cached != nil {
		var results models.Results
		if err := json.Unmarshal(cached, &results); err == nil {
			return &results, nil
		}
	}

	// query db
	tally, err := vote.VoteRepo.GetVoteTally(ctx)
	if err != nil {
		return nil, fmt.Errorf("results: get tally: %w", err)
	}

	total, err := vote.VoteRepo.GetTotalVoteCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("results: get total vote count: %w", err)
	}

	results := &models.Results{
		Tally:      tally,
		TotalVotes: total,
		IsTie:      false,
		Leaders:    leadersFromTally(tally),
		CachedAt:   time.Now(),
	}
	results.IsTie = len(results.Leaders) > 1

	if data, err := json.Marshal(results); err == nil {
		if err := cache.SetResultsCache(ctx, s.rdb, data); err != nil {
			log.Warn().Err(err).Msg("failed to cache results")
		}
	}

	return results, nil
}

func leadersFromTally(tally []models.TallyRow) []models.TallyRow {
	if len(tally) == 0 {
		return nil
	}

	topVoteCount := tally[0].VoteCount
	leaders := make([]models.TallyRow, 0, len(tally))

	for _, row := range tally {
		if row.VoteCount != topVoteCount {
			break
		}
		leaders = append(leaders, row)
	}

	return leaders
}
