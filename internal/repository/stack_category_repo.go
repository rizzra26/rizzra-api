package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/database"
	"github.com/rizzra/api/internal/models"
)

type StackCategoryRepo struct {
	pool *pgxpool.Pool
}

func NewStackCategoryRepo(pool *pgxpool.Pool) *StackCategoryRepo {
	return &StackCategoryRepo{pool: pool}
}

func (r *StackCategoryRepo) List(ctx context.Context) ([]models.StackCategory, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, slug, description, order_index, created_at, updated_at FROM stack_categories WHERE deleted_at IS NULL ORDER BY order_index ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.StackCategory, 0)
	for rows.Next() {
		var c models.StackCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Items = make([]models.StackItem, 0)
		categories = append(categories, c)
	}

	for i := range categories {
		items, err := r.listItems(ctx, categories[i].ID)
		if err != nil {
			return nil, err
		}
		categories[i].Items = items
	}

	return categories, nil
}

func (r *StackCategoryRepo) GetByID(ctx context.Context, id string) (*models.StackCategory, error) {
	var c models.StackCategory
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, description, order_index, created_at, updated_at FROM stack_categories WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	items, err := r.listItems(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Items = items

	return &c, nil
}

func (r *StackCategoryRepo) Create(ctx context.Context, c *models.StackCategory) error {
	if c.Slug == "" {
		c.Slug = database.NormalizeSlug(c.Name)
	}

	var maxOrder int
	r.pool.QueryRow(ctx, `SELECT COALESCE(MAX(order_index), -1) FROM stack_categories WHERE deleted_at IS NULL`).Scan(&maxOrder)
	c.OrderIndex = maxOrder + 1

	return r.pool.QueryRow(ctx,
		`INSERT INTO stack_categories (name, slug, description, order_index) VALUES ($1, $2, $3, $4) RETURNING id, order_index, created_at, updated_at`,
		c.Name, c.Slug, c.Description, c.OrderIndex,
	).Scan(&c.ID, &c.OrderIndex, &c.CreatedAt, &c.UpdatedAt)
}

func (r *StackCategoryRepo) Update(ctx context.Context, id string, c *models.StackCategory) error {
	if c.Slug == "" {
		c.Slug = database.NormalizeSlug(c.Name)
	}

	return r.pool.QueryRow(ctx,
		`UPDATE stack_categories SET name = $1, slug = $2, description = $3, updated_at = NOW() WHERE id = $4 AND deleted_at IS NULL RETURNING order_index, created_at, updated_at`,
		c.Name, c.Slug, c.Description, id,
	).Scan(&c.OrderIndex, &c.CreatedAt, &c.UpdatedAt)
}

func (r *StackCategoryRepo) SoftDelete(ctx context.Context, id string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `UPDATE stack_items SET deleted_at = NOW() WHERE category_id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		// also set deleted_at on category
	}
	tag, err = tx.Exec(ctx, `UPDATE stack_categories SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return tx.Commit(ctx)
}

func (r *StackCategoryRepo) Reorder(ctx context.Context, orders []struct {
	ID         string `json:"id"`
	OrderIndex int    `json:"order_index"`
}) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, o := range orders {
		if _, err := tx.Exec(ctx, `UPDATE stack_categories SET order_index = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`, o.OrderIndex, o.ID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *StackCategoryRepo) listItems(ctx context.Context, categoryID string) ([]models.StackItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, category_id, name, description, order_index, created_at, updated_at FROM stack_items WHERE category_id = $1 AND deleted_at IS NULL ORDER BY order_index ASC`,
		categoryID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.StackItem, 0)
	for rows.Next() {
		var item models.StackItem
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.OrderIndex, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
