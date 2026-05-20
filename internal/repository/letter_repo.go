package repository

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/database"
	"github.com/rizzra/api/internal/models"
)

type LetterRepo struct {
	pool *pgxpool.Pool
}

func NewLetterRepo(pool *pgxpool.Pool) *LetterRepo {
	return &LetterRepo{pool: pool}
}

func (r *LetterRepo) List(ctx context.Context, page, perPage int) ([]models.LetterSummary, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM letters WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	rows, err := r.pool.Query(ctx,
		`SELECT id, slug, title, subtitle, reading_time, created_at, updated_at FROM letters WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	letters := make([]models.LetterSummary, 0)
	for rows.Next() {
		var l models.LetterSummary
		if err := rows.Scan(&l.ID, &l.Slug, &l.Title, &l.Subtitle, &l.ReadingTime, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, 0, err
		}
		letters = append(letters, l)
	}

	return letters, total, nil
}

func (r *LetterRepo) GetByID(ctx context.Context, id string) (*models.Letter, error) {
	var l models.Letter
	err := r.pool.QueryRow(ctx,
		`SELECT id, slug, title, subtitle, content, reading_time, created_at, updated_at, deleted_at FROM letters WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&l.ID, &l.Slug, &l.Title, &l.Subtitle, &l.Content, &l.ReadingTime, &l.CreatedAt, &l.UpdatedAt, &l.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("letter not found: %w", err)
	}
	return &l, nil
}

func (r *LetterRepo) Create(ctx context.Context, l *models.Letter) error {
	l.Slug = r.generateUniqueSlug(ctx, database.NormalizeSlug(l.Title))
	l.ReadingTime = computeReadingTime(l.Content)

	return r.pool.QueryRow(ctx,
		`INSERT INTO letters (slug, title, subtitle, content, reading_time) VALUES ($1, $2, $3, $4, $5) RETURNING id, slug, reading_time, created_at, updated_at`,
		l.Slug, l.Title, l.Subtitle, l.Content, l.ReadingTime,
	).Scan(&l.ID, &l.Slug, &l.ReadingTime, &l.CreatedAt, &l.UpdatedAt)
}

func (r *LetterRepo) Update(ctx context.Context, id string, l *models.Letter) error {
	l.Slug = database.NormalizeSlug(l.Title)
	l.Slug = r.generateUniqueSlug(ctx, l.Slug, id)
	l.ReadingTime = computeReadingTime(l.Content)

	return r.pool.QueryRow(ctx,
		`UPDATE letters SET slug = $1, title = $2, subtitle = $3, content = $4, reading_time = $5, updated_at = NOW() WHERE id = $6 AND deleted_at IS NULL RETURNING id, slug, reading_time, updated_at`,
		l.Slug, l.Title, l.Subtitle, l.Content, l.ReadingTime, id,
	).Scan(&l.ID, &l.Slug, &l.ReadingTime, &l.UpdatedAt)
}

func (r *LetterRepo) SoftDelete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE letters SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("letter not found")
	}
	return nil
}

func (r *LetterRepo) generateUniqueSlug(ctx context.Context, base string, excludeIDs ...string) string {
	slug := base
	if slug == "" {
		slug = "untitled"
	}

	for attempt := 0; attempt < 100; attempt++ {
		candidate := slug
		if attempt > 0 {
			candidate = fmt.Sprintf("%s-%d", slug, attempt+1)
		}

		var query string
		var args []any

		if len(excludeIDs) > 0 {
			query = `SELECT NOT EXISTS(SELECT 1 FROM letters WHERE slug = $1 AND deleted_at IS NULL AND id != $2)`
			args = []any{candidate, excludeIDs[0]}
		} else {
			query = `SELECT NOT EXISTS(SELECT 1 FROM letters WHERE slug = $1 AND deleted_at IS NULL)`
			args = []any{candidate}
		}

		var available bool
		if err := r.pool.QueryRow(ctx, query, args...).Scan(&available); err == nil && available {
			return candidate
		}
	}

	return fmt.Sprintf("%s-%d", slug, rand.Intn(10000))
}

func computeReadingTime(content string) int {
	words := len(strings.Fields(content))
	minutes := math.Ceil(float64(words) / 200.0)
	if minutes < 1 {
		minutes = 1
	}
	return int(minutes)
}
