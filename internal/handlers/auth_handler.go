package handler

import (
	"net/http"
	"user-service/internal/dto"
	"user-service/internal/middleware"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с username и password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Данные для регистрации"
// @Success      201  {object}  map[string]interface{}  "Успешная регистрация"
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user successfully registered",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// Login godoc
// @Summary      Вход пользователя
// @Description  Аутентификация пользователя и выдача JWT токенов
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Данные для входа"
// @Success      200  {object}  dto.TokenResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	h.setAccessTokenCookie(c, tokens.AccessToken)
	h.setRefreshTokenCookie(c, tokens.RefreshToken)

	c.JSON(http.StatusOK, tokens)
}

// Logout godoc
// @Summary      Выход пользователя
// @Description  Отзыв refresh токена пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  string
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "not authorized"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, "successfully logged out")
}

func (h *AuthHandler) setAccessTokenCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("access_token", token, 15*60, "/", "", false, true)
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", token, 7*24*60*60, "/", "", false, true)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}
