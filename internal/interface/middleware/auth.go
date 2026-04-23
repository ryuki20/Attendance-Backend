package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/myuto/attendance-backend/internal/domain/entity"
	"github.com/myuto/attendance-backend/pkg/utils"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization header format",
			})
		}

		token := tokenParts[1]
		claims, err := utils.ValidateJWT(token, m.jwtSecret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		// クレームをコンテキストに保存
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		return next(c)
	}
}

func (m *AuthMiddleware) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != string(entity.RoleAdmin) {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "insufficient permissions",
			})
		}
		return next(c)
	}
}

