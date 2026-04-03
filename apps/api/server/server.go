package server

import (
	"context"
	"net/http"
	"time"

	"github.com/emmanuella-codes/olu/config"
	"github.com/emmanuella-codes/olu/handlers"
	adminhandler "github.com/emmanuella-codes/olu/handlers/admin-handler"
	"github.com/emmanuella-codes/olu/middleware"
	"github.com/emmanuella-codes/olu/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func RunServer(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, rdb *redis.Client, smsSvc *services.SMSService) {
	if cfg.Environment != "development" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type", "X-Request-ID"},
		MaxAge:        5 * time.Minute,
	}))

	candidateSvc := services.NewCandidateService(pool, rdb)
	voteSvc := services.NewVoteService(pool, rdb, smsSvc)
	otpSvc := services.NewOTPService(rdb, smsSvc)
	resultsSvc := services.NewResultsService(pool, rdb)

	candidateHandler := handlers.NewCandidateHandler(candidateSvc)
	otpHandler := handlers.NewOTPHandler(otpSvc, cfg.JWTSecret)
	voteHandler := handlers.NewVoteHandler(voteSvc)
	resultsHandler := handlers.NewResultsHandler(resultsSvc)
	healthHandler := handlers.NewHealthHandler(pool, rdb)
	adminHandler := adminhandler.NewAdminHandler(cfg.AdminJWTSecret)

	r.GET("/health", healthHandler.Health)

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/send-otp", middleware.RateLimit(rdb, "send_otp", 5, time.Minute), otpHandler.Send)
			authGroup.POST("/verify-otp", middleware.RateLimit(rdb, "verify_otp", 10, time.Minute), otpHandler.Verify)
		}

		api.GET("/candidates", candidateHandler.List)
		api.GET("/candidates/:id", candidateHandler.GetByID)
		api.GET("/results", resultsHandler.GetResults)
		api.POST("/vote", middleware.RequireOTPToken(cfg.JWTSecret), voteHandler.Cast)

		adminGroup := api.Group("/admin")
		adminGroup.POST("/login", middleware.RateLimit(rdb, "admin_login", 5, time.Minute), adminHandler.Login)

		adminGroup.Use(middleware.RequireAdminToken(cfg.AdminJWTSecret))
		{
			adminGroup.POST("/admins", adminHandler.CreateAdmin)
			adminGroup.GET("/candidates", adminHandler.AllCandidates)
			adminGroup.POST("/candidates", adminHandler.CreateCandidate)
			adminGroup.PATCH("/candidates/:id", adminHandler.UpdateCandidate)
			adminGroup.DELETE("/candidates/:id", adminHandler.DeactivateCandidate)
			adminGroup.GET("/stats", adminHandler.Stats)
		}
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
