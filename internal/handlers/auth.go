package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizzra/api/internal/repository"
	"github.com/rizzra/api/internal/util"
)

type AuthHandler struct {
	userRepo  *repository.UserRepo
	jwtSecret string
}

func NewAuthHandler(pool *pgxpool.Pool, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userRepo:  repository.NewUserRepo(pool),
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Email    string `json:"email" form:"Email" validate:"required,email"`
	Password string `json:"password" form:"Password" validate:"required,min=6"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	user, err := h.userRepo.FindByEmail(c.Context(), req.Email)
	if err != nil {
		return util.Error(c, 401, "Email or password is incorrect")
	}

	if !util.CheckPassword(req.Password, user.Password) {
		return util.Error(c, 401, "Email or password is incorrect")
	}

	accessToken, accessExp, err := util.GenerateAccessToken(user.ID, h.jwtSecret)
	if err != nil {
		return util.Error(c, 500, "Failed to generate access token")
	}

	refreshToken, _, err := util.GenerateRefreshToken(user.ID, h.jwtSecret)
	if err != nil {
		return util.Error(c, 500, "Failed to generate refresh token")
	}

	return util.OK(c, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp.Unix(),
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" form:"RefreshToken" validate:"required"`
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return util.Error(c, 400, "Invalid request body")
	}

	claims, err := util.ValidateToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		return util.Error(c, 401, "Invalid or expired refresh token")
	}

	accessToken, accessExp, err := util.GenerateAccessToken(claims.Subject, h.jwtSecret)
	if err != nil {
		return util.Error(c, 500, "Failed to generate access token")
	}

	refreshToken, _, err := util.GenerateRefreshToken(claims.Subject, h.jwtSecret)
	if err != nil {
		return util.Error(c, 500, "Failed to generate refresh token")
	}

	return util.OK(c, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    accessExp.Unix(),
	})
}
