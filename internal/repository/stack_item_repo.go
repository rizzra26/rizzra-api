package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
)

type StackItemRepo struct {
	pool *pgxpool.Pool
}

func NewStackItemRepo(pool *pgxpool.Pool) *StackItemRepo {
	return &StackItemRepo{pool: pool}
}

func (r *StackItemRepo) Create(ctx context.Context, item *models.StackItem) error {
	var maxOrder int
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(order_index), -1) FROM stack_items WHERE category_id = $1 AND deleted_at IS NULL`, item.CategoryID,
	).Scan(&maxOrder)
	item.OrderIndex = maxOrder + 1

	return r.pool.QueryRow(ctx,
		`INSERT INTO stack_items (category_id, name, description, order_index) VALUES ($1, $2, $3, $4) RETURNING id, order_index, created_at, updated_at`,
		item.CategoryID, item.Name, item.Description, item.OrderIndex,
	).Scan(&item.ID, &item.OrderIndex, &item.CreatedAt, &item.UpdatedAt)
}

func (r *StackItemRepo) Update(ctx context.Context, id string, item *models.StackItem) error {
	return r.pool.QueryRow(ctx,
		`UPDATE stack_items SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL RETURNING category_id, order_index, created_at, updated_at`,
		item.Name, item.Description, id,
	).Scan(&item.CategoryID, &item.OrderIndex, &item.CreatedAt, &item.UpdatedAt)
}

func (r *StackItemRepo) SoftDelete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE stack_items SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}
