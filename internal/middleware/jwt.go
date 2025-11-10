package middleware

import (
	"net/http"
	"strings"
	"user-service/internal/dto"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthMiddleware(jwtService *services.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		tokenString, err := c.Cookie("access_token")

		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "token is missing",
				})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
					Error: "invalid format Authorization header",
				})
				c.Abort()
				return
			}

			tokenString = parts[1]
		}

		claims, err := m.jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) (int, bool) {
	userID, ok := c.Get("userID")
	if !ok {
		return 0, false
	}
	return userID.(int), true
}

func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, ok := c.Get("username")
	if !ok {
		return "", false
	}

	return username.(string), true
}
