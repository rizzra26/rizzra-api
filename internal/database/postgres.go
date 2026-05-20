package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Connect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("connected to postgres")
	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		content, err := migrationsFS.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
		}

		slog.Info("running migration", "name", entry.Name())
		if _, err := pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", entry.Name(), err)
		}
	}

	slog.Info("migrations complete")
	return nil
}

func SeedAdmin(ctx context.Context, pool *pgxpool.Pool, email, username, password string) error {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		slog.Info("admin user already exists")
		return nil
	}

	_, err = pool.Exec(ctx, `INSERT INTO users (email, username, password, role) VALUES ($1, $2, $3, 'admin')`,
		email, username, password)
	if err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	slog.Info("seeded admin user", "email", email)
	return nil
}

// NormalizeSlug converts a string to a URL-friendly slug.
func NormalizeSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	slug := result.String()
	slug = strings.Join(strings.FieldsFunc(slug, func(r rune) bool { return r == '-' }), "-")
	return strings.Trim(slug, "-")
}
