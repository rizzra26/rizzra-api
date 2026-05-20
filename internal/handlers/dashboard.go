package handlers

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/util"
)

type DashboardHandler struct {
	pool *pgxpool.Pool
}

func NewDashboardHandler(pool *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{pool: pool}
}

type recentLetter struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type monthCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type categoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

func (h *DashboardHandler) Stats(c fiber.Ctx) error {
	ctx := c.Context()

	var lettersCount, projectsCount, categoriesCount, itemsCount int

	if err := h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM letters WHERE deleted_at IS NULL`).Scan(&lettersCount); err != nil {
		return util.Error(c, 500, "Failed to fetch stats")
	}
	if err := h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL`).Scan(&projectsCount); err != nil {
		return util.Error(c, 500, "Failed to fetch stats")
	}
	if err := h.pool.QueryRow(ctx, `SELECT COUNT(*) FROM stack_categories WHERE deleted_at IS NULL`).Scan(&categoriesCount); err != nil {
		return util.Error(c, 500, "Failed to fetch stats")
	}
	err := h.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM stack_items si JOIN stack_categories sc ON si.category_id = sc.id WHERE si.deleted_at IS NULL AND sc.deleted_at IS NULL`,
	).Scan(&itemsCount)
	if err != nil {
		return util.Error(c, 500, "Failed to fetch stats")
	}

	recentLetters := make([]recentLetter, 0)
	rows, err := h.pool.Query(ctx, `SELECT id, title, created_at FROM letters WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT 5`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rl recentLetter
			if err := rows.Scan(&rl.ID, &rl.Title, &rl.CreatedAt); err == nil {
				recentLetters = append(recentLetters, rl)
			}
		}
	}

	lettersByMonth := make([]monthCount, 0)
	lrows, err := h.pool.Query(ctx,
		`SELECT to_char(date_trunc('month', created_at), 'YYYY-MM') AS month, COUNT(*) FROM letters WHERE deleted_at IS NULL GROUP BY date_trunc('month', created_at) ORDER BY date_trunc('month', created_at) ASC`)
	if err == nil {
		defer lrows.Close()
		for lrows.Next() {
			var mc monthCount
			if err := lrows.Scan(&mc.Month, &mc.Count); err == nil {
				lettersByMonth = append(lettersByMonth, mc)
			}
		}
	}

	projectsByCategory := make([]categoryCount, 0)
	prows, err := h.pool.Query(ctx,
		`SELECT category, COUNT(*) FROM projects WHERE deleted_at IS NULL GROUP BY category`)
	if err == nil {
		defer prows.Close()
		for prows.Next() {
			var cc categoryCount
			if err := prows.Scan(&cc.Category, &cc.Count); err == nil {
				projectsByCategory = append(projectsByCategory, cc)
			}
		}
	}

	return util.OK(c, fiber.Map{
		"letters_count":         lettersCount,
		"projects_count":        projectsCount,
		"stack_categories_count": categoriesCount,
		"stack_items_count":     itemsCount,
		"recent_letters":        recentLetters,
		"letters_by_month":      lettersByMonth,
		"projects_by_category":  projectsByCategory,
	})
}
