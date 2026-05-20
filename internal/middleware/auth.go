package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rizzra/api/internal/util"
)

func Auth(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return util.Error(c, 401, "Missing authorization header")
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return util.Error(c, 401, "Invalid authorization format")
		}

		claims, err := util.ValidateToken(parts[1], secret)
		if err != nil {
			return util.Error(c, 401, "Invalid or expired token")
		}

		c.Locals("userID", claims.Subject)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}
