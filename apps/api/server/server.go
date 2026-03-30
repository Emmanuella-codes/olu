package server

import (
	"context"
	"net/http"
	"time"

	"github.com/emmanuella-codes/olu/config"
	"github.com/emmanuella-codes/olu/handlers"
	"github.com/emmanuella-codes/olu/middleware"
	"github.com/emmanuella-codes/olu/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func RunServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, rdb *redis.Client) {
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		// AllowAllOrigins: ,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           5 * time.Minute,
	}))

	candidateSvc := services.NewCandidateService(pool, rdb)
	// voteSvc := services.NewVoteService(pool, rdb)
	resultsSvc := services.NewResultsService(pool, rdb)

	candidateHandler := handlers.NewCandidateHandler(candidateSvc)
	resultsHandler := handlers.NewResultsHandler(resultsSvc)
	healthHandler := handlers.NewHealthHandler(pool, rdb)

	r.GET("/health", healthHandler.Health)

	api := r.Group("/api/v1")
	{
		api.GET("/candidates", candidateHandler.List)
		api.GET("/candidates/:id", candidateHandler.GetByID)
		api.GET("/results", resultsHandler.GetResults)

		// admin
		// adminGroup := api.Group("/admin")
		// {
		// 	adminGroup.POST("/login")
		// }
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown server")
		}
		log.Info().Msg("server shutdown complete")
	}()

	log.Info().Str("port", cfg.Port).Str("env", cfg.Environment).Msg("server starting")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server error")
	}
	log.Info().Msg("server stopped")
}
