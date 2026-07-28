package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerRequest struct {
	Email       string `json:"email" example:"student@test.com"`
	Password    string `json:"password" example:"password1"`
	DisplayName string `json:"displayName" example:"Студент"`
}

type loginRequest struct {
	Email    string `json:"email" example:"student@test.com"`
	Password string `json:"password" example:"password1"`
}

// Register регистрация пользователя.
//
// @Summary     Регистрация
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body registerRequest true "Данные регистрации"
// @Success     201 {object} domain.AuthResponse
// @Failure     409 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	resp, err := h.svc.Register(c.Request.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// Login вход в систему.
//
// @Summary     Вход
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body loginRequest true "Учётные данные"
// @Success     200 {object} domain.AuthResponse
// @Failure     401 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	resp, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Me текущий пользователь.
//
// @Summary     Текущий пользователь
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} domain.User
// @Failure     401 {object} ErrorResponse
// @Router      /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.svc.Me(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
