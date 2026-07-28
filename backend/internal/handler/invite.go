package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type InviteHandler struct {
	svc *service.InviteService
}

func NewInviteHandler(svc *service.InviteService) *InviteHandler {
	return &InviteHandler{svc: svc}
}

// Create создание приглашения в проект.
//
// @Summary     Создать приглашение
// @Tags        invites
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     201 {object} InviteResponse
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Router      /projects/{projectId}/invites [post]
func (h *InviteHandler) Create(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	invite, err := h.svc.Create(c.Request.Context(), projectID, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, invite)
}

// Accept принятие приглашения.
//
// @Summary     Принять приглашение
// @Tags        invites
// @Produce     json
// @Security    BearerAuth
// @Param       token path string true "Токен приглашения"
// @Success     200 {object} domain.Project
// @Failure     401 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /invites/{token}/accept [post]
func (h *InviteHandler) Accept(c *gin.Context) {
	token := c.Param("token")
	project, err := h.svc.Accept(c.Request.Context(), token, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

// ссылка на domain.Project для swag
var _ domain.Project
