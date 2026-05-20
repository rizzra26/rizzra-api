package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type LetterHandler struct {
	repo *repository.LetterRepo
}

func NewLetterHandler(pool *pgxpool.Pool) *LetterHandler {
	return &LetterHandler{repo: repository.NewLetterRepo(pool)}
}

func (h *LetterHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	letters, total, err := h.repo.List(c.Context(), page, perPage)
	if err != nil {
		return util.Error(c, 500, "Failed to fetch letters")
	}

	totalPages := (total + perPage - 1) / perPage
	return util.PaginatedOK(c, letters, util.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *LetterHandler) Get(c fiber.Ctx) error {
	letter, err := h.repo.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return util.Error(c, 404, "Letter not found")
	}
	return util.OK(c, letter)
}

type createLetterRequest struct {
	Title    string `json:"title" validate:"required,min=1"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content" validate:"required"`
}

func (h *LetterHandler) Create(c fiber.Ctx) error {
	var req createLetterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	letter := &models.Letter{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		Content:  req.Content,
	}
	if err := h.repo.Create(c.Context(), letter); err != nil {
		return util.Error(c, 500, "Failed to create letter")
	}

	return util.Created(c, letter)
}

func (h *LetterHandler) Update(c fiber.Ctx) error {
	var req createLetterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	letter := &models.Letter{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		Content:  req.Content,
	}
	if err := h.repo.Update(c.Context(), c.Params("id"), letter); err != nil {
		return util.Error(c, 500, "Failed to update letter")
	}

	return util.OK(c, letter)
}

func (h *LetterHandler) Delete(c fiber.Ctx) error {
	if err := h.repo.SoftDelete(c.Context(), c.Params("id")); err != nil {
		return util.Error(c, 404, "Letter not found")
	}
	return util.Deleted(c)
}
