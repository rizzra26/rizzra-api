package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type StackHandler struct {
	itemRepo *repository.StackItemRepo
}

func NewStackHandler(pool *pgxpool.Pool) *StackHandler {
	return &StackHandler{
		itemRepo: repository.NewStackItemRepo(pool),
	}
}

func (h *StackHandler) ListItems(c fiber.Ctx) error {
	items, err := h.itemRepo.List(c.Context())
	if err != nil {
		return util.Error(c, 500, "Failed to fetch stack items")
	}
	return util.OK(c, items)
}

type createItemRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

func (h *StackHandler) CreateItem(c fiber.Ctx) error {
	var req createItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	item := &models.StackItem{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.itemRepo.Create(c.Context(), item); err != nil {
		return util.Error(c, 500, "Failed to create item")
	}

	return util.Created(c, item)
}

type updateItemRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

func (h *StackHandler) UpdateItem(c fiber.Ctx) error {
	var req updateItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	item := &models.StackItem{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.itemRepo.Update(c.Context(), c.Params("id"), item); err != nil {
		return util.Error(c, 500, "Failed to update item")
	}

	return util.OK(c, item)
}

func (h *StackHandler) DeleteItem(c fiber.Ctx) error {
	if err := h.itemRepo.SoftDelete(c.Context(), c.Params("id")); err != nil {
		return util.Error(c, 404, "Stack item not found")
	}
	return util.Deleted(c)
}
