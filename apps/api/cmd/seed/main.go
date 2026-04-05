package main

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/emmanuella-codes/olu/config"
	"github.com/emmanuella-codes/olu/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	email := flag.String("email", "", "admin email address (required)")
	password := flag.String("password", "", "admin password (required)")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if *email == "" || *password == "" {
		log.Fatal().Msg("usage: go run ./cmd/seed -email <email> -password <password>")
	}

	*email = strings.ToLower(strings.TrimSpace(*email))
	if len(*password) < 8 {
		log.Fatal().Msg("password must be at least 8 characters")
	}

	cfg, err := config.LoadForMigration()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM admins WHERE email = $1)", *email).Scan(&exists)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to check existing admin")
	}
	if exists {
		log.Fatal().Str("email", *email).Msg("admin with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to hash password")
	}

	var id uuid.UUID
	err = pool.QueryRow(ctx, `
		INSERT INTO admins (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`, uuid.New(), *email, string(hash)).Scan(&id)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create admin")
	}

	log.Info().Str("id", id.String()).Str("email", *email).Msg("admin created successfully")
}
