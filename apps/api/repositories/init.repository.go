package repositories

import (
	"log"

	"github.com/emmanuella-codes/olu/repositories/admin"
	"github.com/emmanuella-codes/olu/repositories/vote"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRepository(pool *pgxpool.Pool, logger *log.Logger) {
	vote.InitVoteRepo(pool, logger)
	admin.InitAdminRepo(pool, logger)
}
