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

func (r *StackItemRepo) List(ctx context.Context) ([]models.StackItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, order_index, created_at, updated_at FROM stack_items WHERE deleted_at IS NULL ORDER BY order_index ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.StackItem, 0)
	for rows.Next() {
		var item models.StackItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.OrderIndex, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *StackItemRepo) Create(ctx context.Context, item *models.StackItem) error {
	var maxOrder int
	r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(order_index), -1) FROM stack_items WHERE deleted_at IS NULL`,
	).Scan(&maxOrder)
	item.OrderIndex = maxOrder + 1

	return r.pool.QueryRow(ctx,
		`INSERT INTO stack_items (name, description, order_index) VALUES ($1, $2, $3) RETURNING id, order_index, created_at, updated_at`,
		item.Name, item.Description, item.OrderIndex,
	).Scan(&item.ID, &item.OrderIndex, &item.CreatedAt, &item.UpdatedAt)
}

func (r *StackItemRepo) Update(ctx context.Context, id string, item *models.StackItem) error {
	return r.pool.QueryRow(ctx,
		`UPDATE stack_items SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND deleted_at IS NULL RETURNING order_index, created_at, updated_at`,
		item.Name, item.Description, id,
	).Scan(&item.OrderIndex, &item.CreatedAt, &item.UpdatedAt)
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
