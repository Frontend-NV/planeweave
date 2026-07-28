package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/middleware"
	"github.com/planeweave/planeweave/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type createTaskRequest struct {
	Title      string  `json:"title" example:"Задача A"`
	Status     string  `json:"status" example:"todo" enums(todo,in_progress,done)`
	AssigneeID *string `json:"assigneeId"`
	X          float64 `json:"x" example:"100"`
	Y          float64 `json:"y" example:"100"`
}

type patchTaskRequest struct {
	Title      *string  `json:"title"`
	Status     *string  `json:"status" enums(todo,in_progress,done)`
	AssigneeID *string  `json:"assigneeId"`
	X          *float64 `json:"x"`
	Y          *float64 `json:"y"`
}

// List список задач проекта.
//
// @Summary     Список задач
// @Tags        tasks
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Success     200 {array} domain.Task
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Router      /projects/{projectId}/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	tasks, err := h.svc.List(c.Request.Context(), projectID, middleware.UserID(c))
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Create создание задачи.
//
// @Summary     Создать задачу
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       projectId path string true "ID проекта"
// @Param       body body createTaskRequest true "Задача"
// @Success     201 {object} domain.Task
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     422 {object} ErrorResponse
// @Router      /projects/{projectId}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	projectID, ok := parseProjectID(c)
	if !ok {
		return
	}
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	task, err := h.svc.Create(c.Request.Context(), projectID, middleware.UserID(c), service.CreateTaskInput{
		Title:      req.Title,
		Status:     req.Status,
		AssigneeID: req.AssigneeID,
		X:          req.X,
		Y:          req.Y,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, task)
}

// Patch изменение задачи.
//
// @Summary     Изменить задачу
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       taskId path string true "ID задачи"
// @Param       body body patchTaskRequest true "Поля для изменения"
// @Success     200 {object} domain.Task
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /tasks/{taskId} [patch]
func (h *TaskHandler) Patch(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}
	var req patchTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, domain.ErrInvalidRequest)
		return
	}
	task, _, err := h.svc.Update(c.Request.Context(), taskID, middleware.UserID(c), domain.TaskPatch{
		Title:      req.Title,
		Status:     req.Status,
		AssigneeID: req.AssigneeID,
		X:          req.X,
		Y:          req.Y,
	})
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

// Delete удаление задачи.
//
// @Summary     Удалить задачу
// @Tags        tasks
// @Security    BearerAuth
// @Param       taskId path string true "ID задачи"
// @Success     204
// @Failure     401 {object} ErrorResponse
// @Failure     403 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Router      /tasks/{taskId} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, ok := parseTaskID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), taskID, middleware.UserID(c)); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseTaskID(c *gin.Context) (uuid.UUID, bool) {
	id, ok := parseUUID(c, "taskId")
	if !ok {
		return uuid.Nil, false
	}
	parsed, _ := uuid.Parse(id)
	return parsed, true
}
