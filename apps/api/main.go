package main

import (
	"context"
	"flag"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emmanuella-codes/olu/cache"
	"github.com/emmanuella-codes/olu/config"
	"github.com/emmanuella-codes/olu/db"
	"github.com/emmanuella-codes/olu/migrations"
	"github.com/emmanuella-codes/olu/repositories"
	"github.com/emmanuella-codes/olu/server"
	"github.com/emmanuella-codes/olu/services"
	"github.com/emmanuella-codes/olu/sms"
	"github.com/emmanuella-codes/olu/workers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	migrateFlag := flag.String("migrate", "", "run database migration: up or down")
	forceDownFlag := flag.Bool("force-down", false, "allow destructive migration down")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *migrateFlag != "" {
		cfg, err := config.LoadForMigration()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load migration config")
		}

		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer cancel()

		pool, err := db.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}
		defer pool.Close()

		switch *migrateFlag {
		case "up":
			if err := migrations.Up(ctx, pool); err != nil {
				log.Fatal().Err(err).Msg("migration up failed")
			}
			log.Info().Msg("migration up complete")
		case "down":
			if !*forceDownFlag {
				log.Fatal().Msg("migration down is destructive; rerun with -force-down")
			}
			if err := migrations.Down(ctx, pool); err != nil {
				log.Fatal().Err(err).Msg("migration down failed")
			}
			log.Info().Msg("migration down complete")
		default:
			log.Fatal().Str("value", *migrateFlag).Msg("invalid -migrate value: use 'up' or 'down'")
		}
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if cfg.Environment == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	rdb, err := cache.NewRedis(cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close redis connection")
		}
	}()

	repositories.InitRepository(pool, stdlog.New(os.Stderr, "", stdlog.LstdFlags))

	smsProvider, err := sms.Build(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to build sms provider")
	}
	smsSvc := services.NewSMSService(smsProvider, cfg.SMSFrom)
	smsWorker := workers.NewSMSWorker(smsSvc)
	go smsWorker.RunQueueWorker(ctx, 20, 2*time.Second)

	server.RunServer(ctx, cfg, pool, rdb, smsSvc)
}
