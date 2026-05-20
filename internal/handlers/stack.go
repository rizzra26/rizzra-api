package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type StackHandler struct {
	catRepo  *repository.StackCategoryRepo
	itemRepo *repository.StackItemRepo
}

func NewStackHandler(pool *pgxpool.Pool) *StackHandler {
	return &StackHandler{
		catRepo:  repository.NewStackCategoryRepo(pool),
		itemRepo: repository.NewStackItemRepo(pool),
	}
}

// ---- Categories ----

func (h *StackHandler) ListCategories(c fiber.Ctx) error {
	categories, err := h.catRepo.List(c.Context())
	if err != nil {
		return util.Error(c, 500, "Failed to fetch categories")
	}
	return util.OK(c, categories)
}

type createCategoryRequest struct {
	Name        string  `json:"name" validate:"required"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
}

func (h *StackHandler) CreateCategory(c fiber.Ctx) error {
	var req createCategoryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	cat := &models.StackCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}
	if err := h.catRepo.Create(c.Context(), cat); err != nil {
		return util.Error(c, 500, "Failed to create category")
	}

	return util.Created(c, cat)
}

func (h *StackHandler) UpdateCategory(c fiber.Ctx) error {
	var req createCategoryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	cat := &models.StackCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}
	if err := h.catRepo.Update(c.Context(), c.Params("id"), cat); err != nil {
		return util.Error(c, 500, "Failed to update category")
	}

	return util.OK(c, cat)
}

func (h *StackHandler) DeleteCategory(c fiber.Ctx) error {
	if err := h.catRepo.SoftDelete(c.Context(), c.Params("id")); err != nil {
		return util.Error(c, 404, "Category not found")
	}
	return util.Deleted(c)
}

// ---- Items ----

type createItemRequest struct {
	CategoryID  string  `json:"category_id" validate:"required,uuid"`
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
}

func (h *StackHandler) CreateItem(c fiber.Ctx) error {
	var req createItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	item := &models.StackItem{
		CategoryID:  req.CategoryID,
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
