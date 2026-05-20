package util

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c fiber.Ctx, data interface{}) error {
	return c.Status(200).JSON(APIResponse{Data: data})
}

func Created(c fiber.Ctx, data interface{}) error {
	return c.Status(201).JSON(APIResponse{Data: data})
}

func NoContent(c fiber.Ctx) error {
	return c.Status(204).Send(nil)
}

func Deleted(c fiber.Ctx) error {
	return c.Status(200).JSON(APIResponse{Message: "Deleted successfully"})
}

func Error(c fiber.Ctx, status int, errStr string) error {
	return c.Status(status).JSON(APIResponse{Error: errStr})
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func PaginatedOK(c fiber.Ctx, data interface{}, meta PaginationMeta) error {
	return c.Status(200).JSON(APIResponse{Data: data, Meta: meta})
}
