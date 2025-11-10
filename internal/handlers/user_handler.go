package handler

import (
	"net/http"
	"strconv"
	"user-service/internal/dto"
	"user-service/internal/middleware"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetStatus godoc
// @Summary      Получить статус пользователя
// @Description  Возвращает информацию о пользователе и его выполненных заданиях
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  dto.UserStatusResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /api/users/{id}/status [get]
func (h *UserHandler) GetStatus(c *gin.Context) {
	currentUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "not authorized"})
		return
	}

	idParam := c.Param("id")
	requestedUserID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid ID"})
		return
	}

	if currentUserID != requestedUserID {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		return
	}

	response, err := h.userService.GetUserInfoByID(c.Request.Context(), requestedUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLeaderboard godoc
// @Summary      Получить топ пользователей
// @Description  Возвращает список пользователей отсортированных по балансу
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "Количество пользователей"  default(10)
// @Security     BearerAuth
// @Success      200  {object}  string
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/users/leaderboard [get]
func (h *UserHandler) GetLeaderBoard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	response, err := h.userService.GetLeaderBoard(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetLeaderboard godoc
// @Summary      Выполнить задание, получить поинты
// @Description  Возвращает сообщение при успешном выполнении
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "task ID"
// @Security     BearerAuth
// @Success      200  {object}  string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/users/{id}/task/complete [post]
func (h *UserHandler) CompleteTask(c *gin.Context) {
	currentID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "not authorized"})
		return
	}

	idParam := c.Param("id")
	taskIDParam, _ := strconv.Atoi(idParam)

	if err := h.userService.CompleteTask(c.Request.Context(), currentID, taskIDParam); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "task completed")
}

// AddReferrer godoc
// @Summary      Ввести реферальный код
// @Description  Устанавливает реферера для пользователя и начисляет бонусы
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/users/{id}/referrer [post]
func (h *UserHandler) AddReferrer(c *gin.Context) {
	currentUserID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "not authorized"})
		return
	}

	idParam := c.Param("id")
	requestedUserID, _ := strconv.Atoi(idParam)

	if err := h.userService.AddReferrer(c.Request.Context(), currentUserID, requestedUserID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "referrer added")
}
