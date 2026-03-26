package repository

import (
	"log"

	"github.com/emmanuella-codes/olu/repository/vote"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRepository(pool *pgxpool.Pool, logger *log.Logger) {
	vote.InitVoteRepo(pool, logger)
}
