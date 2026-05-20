package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type UploadHandler struct {
	projectRepo *repository.ProjectRepo
	uploadDir   string
}

func NewUploadHandler(projectRepo *repository.ProjectRepo, uploadDir string) *UploadHandler {
	return &UploadHandler{projectRepo: projectRepo, uploadDir: uploadDir}
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

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return util.Error(c, 400, "Invalid file type. Allowed: png, jpg, jpeg, webp")
	}

	if file.Size > 5*1024*1024 {
		return util.Error(c, 400, "File too large. Max 5MB")
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(h.uploadDir, filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return util.Error(c, 500, "Failed to save file")
	}

	url := "/uploads/" + filename

	projectID := c.FormValue("project_id")
	if projectID != "" {
		if err := h.projectRepo.SetCoverURL(c.Context(), projectID, url); err != nil {
			return util.Error(c, 500, "Failed to update project cover")
		}
	}

	return util.OK(c, fiber.Map{
		"url":      url,
		"filename": filename,
	})
}
