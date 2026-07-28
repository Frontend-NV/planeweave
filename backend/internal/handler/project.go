package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

type titleRequest struct {
	Title string `json:"title" example:"Мой проект"`
}

// List список проектов пользователя.
//
// @Summary     Список проектов
// @Tags        projects
// @Produce     json
// @Security    BearerAuth
// @Success     200 {array} domain.Project
// @Failure     401 {object} ErrorResponse
// @Router      /projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// Create создание проекта.
//
// @Summary     Создать проект
// @Tags        projects
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body titleRequest true "Название"
// @Success     201 {object} domain.Project
// @Failure     401 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	var req titleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	project, err := h.svc.Create(c.Request.Context(), middleware.UserID(c), req.Title)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, project)
}

// Get получение проекта.
//
// @Summary     Получить проект
// @Tags        projects
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     200 {object} domain.Project
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /projects/{projectId} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	project, err := h.svc.Get(c.Request.Context(), projectID, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

// Patch переименование проекта.
//
// @Summary     Переименовать проект
// @Tags        projects
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Param       body body titleRequest true "Название"
// @Success     200 {object} domain.Project
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /projects/{projectId} [patch]
func (h *ProjectHandler) Patch(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	var req titleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	project, err := h.svc.UpdateTitle(c.Request.Context(), projectID, middleware.UserID(c), req.Title)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, project)
}

// Delete удаление проекта (только владелец).
//
// @Summary     Удалить проект
// @Tags        projects
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     204
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /projects/{projectId} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), projectID, middleware.UserID(c)); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Members список участников проекта.
//
// @Summary     Участники проекта
// @Tags        projects
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     200 {array} domain.ProjectMember
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Router      /projects/{projectId}/members [get]
func (h *ProjectHandler) Members(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	members, err := h.svc.Members(c.Request.Context(), projectID, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}

func parseProjectID(c *gin.Context) (uuid.UUID, bool) {
	id, ok := parseUUID(c, "projectId")
	if !ok {
		return uuid.Nil, false
	}
	parsed, _ := uuid.Parse(id)
	return parsed, true
}
