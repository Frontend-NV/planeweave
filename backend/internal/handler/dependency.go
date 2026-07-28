package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type DependencyHandler struct {
	svc *service.DependencyService
}

func NewDependencyHandler(svc *service.DependencyService) *DependencyHandler {
	return &DependencyHandler{svc: svc}
}

type createDependencyRequest struct {
	FromTaskID string `json:"fromTaskId" example:"uuid-from"`
	ToTaskID   string `json:"toTaskId" example:"uuid-to"`
}

// List список зависимостей проекта.
//
// @Summary     Список зависимостей
// @Tags        dependencies
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     200 {array} domain.Dependency
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Router      /projects/{projectId}/dependencies [get]
func (h *DependencyHandler) List(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	items, err := h.svc.List(c.Request.Context(), projectID, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

// Create создание зависимости между задачами.
//
// @Summary     Создать зависимость
// @Tags        dependencies
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Param       body body createDependencyRequest true "from блокирует to"
// @Success     201 {object} domain.Dependency
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     409 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /projects/{projectId}/dependencies [post]
func (h *DependencyHandler) Create(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	var req createDependencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	dep, err := h.svc.Create(c.Request.Context(), projectID, middleware.UserID(c), req.FromTaskID, req.ToTaskID)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dep)
}

// Delete удаление зависимости.
//
// @Summary     Удалить зависимость
// @Tags        dependencies
// @Security    BearerAuth
// @Param       dependencyId path string true "ID зависимости"
// @Success     204
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /dependencies/{dependencyId} [delete]
func (h *DependencyHandler) Delete(c *gin.Context) {
	id, ok := parseUUID(c, "dependencyId")
	if !ok {
		return
	}
	dependencyID, _ := uuid.Parse(id)
	if err := h.svc.Delete(c.Request.Context(), dependencyID, middleware.UserID(c)); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
