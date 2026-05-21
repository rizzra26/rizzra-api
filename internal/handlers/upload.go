package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizzra/api/internal/cloudinary"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type UploadHandler struct {
	projectRepo *repository.ProjectRepo
	cld         *cloudinary.Service
}

func NewUploadHandler(projectRepo *repository.ProjectRepo, cld *cloudinary.Service) *UploadHandler {
	return &UploadHandler{projectRepo: projectRepo, cld: cld}
}

var allowedExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

func (h *UploadHandler) Cover(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return util.Error(c, 400, "No file uploaded")
	}

	ext := ""
	for allowed := range allowedExts {
		if strings.HasSuffix(strings.ToLower(file.Filename), allowed) {
			ext = allowed
			break
		}
	}
	if ext == "" {
		return util.Error(c, 400, "Invalid file type. Allowed: png, jpg, jpeg, webp")
	}

	if file.Size > 5*1024*1024 {
		return util.Error(c, 400, "File too large. Max 5MB")
	}

	fd, err := file.Open()
	if err != nil {
		return util.Error(c, 500, "Failed to open file")
	}
	defer fd.Close()

	url, err := h.cld.Upload(c.Context(), fd, file.Filename)
	if err != nil {
		return util.Error(c, 500, "Failed to upload file")
	}

	projectID := c.FormValue("project_id")
	if projectID != "" {
		if err := h.projectRepo.SetCoverURL(c.Context(), projectID, url); err != nil {
			return util.Error(c, 500, "Failed to update project cover")
		}
	}

	return util.OK(c, fiber.Map{
		"url": url,
	})
}
