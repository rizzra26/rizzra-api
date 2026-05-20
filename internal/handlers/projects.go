package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/models"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type ProjectHandler struct {
	repo *repository.ProjectRepo
}

func NewProjectHandler(pool *pgxpool.Pool) *ProjectHandler {
	return &ProjectHandler{repo: repository.NewProjectRepo(pool)}
}

func (h *ProjectHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	projects, total, err := h.repo.List(c.Context(), page, perPage)
	if err != nil {
		return util.Error(c, 500, "Failed to fetch projects")
	}

	totalPages := (total + perPage - 1) / perPage
	return util.PaginatedOK(c, projects, util.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	})
}

func (h *ProjectHandler) Get(c fiber.Ctx) error {
	project, err := h.repo.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		return util.Error(c, 404, "Project not found")
	}
	return util.OK(c, project)
}

type createProjectRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Tech        []string `json:"tech"`
	GithubURL   *string  `json:"github_url"`
	DemoURL     *string  `json:"demo_url"`
	CoverURL    *string  `json:"cover_url"`
	Category    string   `json:"category" validate:"oneof=ai pure"`
}

func (h *ProjectHandler) Create(c fiber.Ctx) error {
	var req createProjectRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}
	if req.Category == "" {
		req.Category = "pure"
	}
	if req.Tech == nil {
		req.Tech = []string{}
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Tech:        req.Tech,
		GithubURL:   req.GithubURL,
		DemoURL:     req.DemoURL,
		CoverURL:    req.CoverURL,
	}
	if err := h.repo.Create(c.Context(), project); err != nil {
		return util.Error(c, 500, "Failed to create project")
	}

	return util.Created(c, project)
}

func (h *ProjectHandler) Update(c fiber.Ctx) error {
	var req createProjectRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}
	if req.Category == "" {
		req.Category = "pure"
	}
	if req.Tech == nil {
		req.Tech = []string{}
	}

	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Tech:        req.Tech,
		GithubURL:   req.GithubURL,
		DemoURL:     req.DemoURL,
		CoverURL:    req.CoverURL,
	}
	if err := h.repo.Update(c.Context(), c.Params("id"), project); err != nil {
		return util.Error(c, 500, "Failed to update project")
	}

	return util.OK(c, project)
}

func (h *ProjectHandler) Delete(c fiber.Ctx) error {
	if err := h.repo.SoftDelete(c.Context(), c.Params("id")); err != nil {
		return util.Error(c, 404, "Project not found")
	}
	return util.Deleted(c)
}

type reorderRequest struct {
	Order []models.ProjectReorderItem `json:"order" validate:"required,min=1"`
}

func (h *ProjectHandler) Reorder(c fiber.Ctx) error {
	var req reorderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	if err := h.repo.Reorder(c.Context(), req.Order); err != nil {
		return util.Error(c, 500, "Failed to reorder projects")
	}

	return util.OK(c, fiber.Map{"message": "Reordered successfully"})
}
