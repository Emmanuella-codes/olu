package migrations

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed *.sql
var migrationFiles embed.FS

const migrationsTable = "schema_migrations"

func Up(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ensureMigrationsTable(ctx, pool); err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, pool)
	if err != nil {
		return err
	}

	files, err := orderedFiles("up")
	if err != nil {
		return err
	}

	for _, file := range files {
		version, err := migrationVersion(file)
		if err != nil {
			return err
		}
		if applied[version] {
			continue
		}

		sql, err := fs.ReadFile(migrationFiles, file)
		if err != nil {
			return fmt.Errorf("migrate up: read %s: %w", file, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("migrate up: begin tx: %w", err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migrate up: apply %s: %w", file, err)
		}
		if _, err := tx.Exec(ctx, "INSERT INTO "+migrationsTable+" (version) VALUES ($1)", version); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migrate up: record %s: %w", file, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("migrate up: commit %s: %w", file, err)
		}
	}

	return nil
}

func Down(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ensureMigrationsTable(ctx, pool); err != nil {
		return err
	}

	versions, err := appliedVersionsOrdered(ctx, pool)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		return nil
	}

	version := versions[len(versions)-1]
	file := version + ".down.sql"
	sql, err := fs.ReadFile(migrationFiles, file)
	if err != nil {
		return fmt.Errorf("migrate down: read %s: %w", file, err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("migrate down: begin tx: %w", err)
	}
	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("migrate down: apply %s: %w", file, err)
	}
	if _, err := tx.Exec(ctx, "DELETE FROM "+migrationsTable+" WHERE version = $1", version); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("migrate down: delete version %s: %w", version, err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("migrate down: commit %s: %w", file, err)
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}
	return nil
}

func appliedVersions(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	rows, err := pool.Query(ctx, "SELECT version FROM "+migrationsTable)
	if err != nil {
		return nil, fmt.Errorf("query applied versions: %w", err)
	}
	defer rows.Close()

	out := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied version: %w", err)
		}
		out[version] = true
	}
	return out, rows.Err()
}

func appliedVersionsOrdered(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	rows, err := pool.Query(ctx, "SELECT version FROM "+migrationsTable+" ORDER BY version ASC")
	if err != nil {
		return nil, fmt.Errorf("query applied versions ordered: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied version: %w", err)
		}
		out = append(out, version)
	}
	return out, rows.Err()
}

func orderedFiles(direction string) ([]string, error) {
	entries, err := fs.Glob(migrationFiles, "*."+direction+".sql")
	if err != nil {
		return nil, fmt.Errorf("glob %s migrations: %w", direction, err)
	}
	sort.Strings(entries)
	return entries, nil
}

func migrationVersion(file string) (string, error) {
	for _, suffix := range []string{".up.sql", ".down.sql"} {
		if strings.HasSuffix(file, suffix) {
			version := strings.TrimSuffix(file, suffix)
			if version == "" {
				return "", fmt.Errorf("invalid migration filename %q", file)
			}
			return version, nil
		}
	}
	return "", fmt.Errorf("invalid migration filename %q", file)
}
