package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
)

type ProjectRepo struct {
	pool *pgxpool.Pool
}

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{pool: pool}
}

func (r *ProjectRepo) List(ctx context.Context, page, perPage int) ([]models.Project, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}

	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, tech, github_url, demo_url, cover_url, order_index, created_at, updated_at FROM projects WHERE deleted_at IS NULL ORDER BY order_index ASC, created_at DESC LIMIT $1 OFFSET $2`,
		perPage, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Tech, &p.GithubURL, &p.DemoURL, &p.CoverURL, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}

	return projects, total, nil
}

func (r *ProjectRepo) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var p models.Project
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, tech, github_url, demo_url, cover_url, order_index, created_at, updated_at, deleted_at FROM projects WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Tech, &p.GithubURL, &p.DemoURL, &p.CoverURL, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	return &p, nil
}

func (r *ProjectRepo) Create(ctx context.Context, p *models.Project) error {
	// Get the max order_index
	var maxOrder int
	r.pool.QueryRow(ctx, `SELECT COALESCE(MAX(order_index), -1) FROM projects WHERE deleted_at IS NULL`).Scan(&maxOrder)
	p.OrderIndex = maxOrder + 1

	return r.pool.QueryRow(ctx,
		`INSERT INTO projects (name, description, tech, github_url, demo_url, cover_url, category, order_index) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, order_index, created_at, updated_at`,
		p.Name, p.Description, p.Tech, p.GithubURL, p.DemoURL, p.CoverURL, p.OrderIndex,
	).Scan(&p.ID, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepo) Update(ctx context.Context, id string, p *models.Project) error {
	return r.pool.QueryRow(ctx,
		`UPDATE projects SET name = $1, description = $2, tech = $3, github_url = $4, demo_url = $5, cover_url = $6, category = $7, updated_at = NOW() WHERE id = $8 AND deleted_at IS NULL RETURNING id, order_index, created_at, updated_at`,
		p.Name, p.Description, p.Tech, p.GithubURL, p.DemoURL, p.CoverURL, id,
	).Scan(&p.ID, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepo) SoftDelete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE projects SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (r *ProjectRepo) Reorder(ctx context.Context, items []models.ProjectReorderItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, item := range items {
		if _, err := tx.Exec(ctx, `UPDATE projects SET order_index = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, item.OrderIndex, item.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *ProjectRepo) SetCoverURL(ctx context.Context, id, url string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE projects SET cover_url = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, url, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}
