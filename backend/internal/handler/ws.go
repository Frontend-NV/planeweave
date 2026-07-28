package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/planeweave/planeweave/internal/domain"
	"github.com/planeweave/planeweave/internal/pkg/jwt"
	"github.com/planeweave/planeweave/internal/service"
	"github.com/planeweave/planeweave/internal/ws"
)

type WSHandler struct {
	tokens  *jwt.Service
	projects *service.ProjectService
	hub     *ws.Hub
}

func NewWSHandler(tokens *jwt.Service, projects *service.ProjectService, hub *ws.Hub) *WSHandler {
	return &WSHandler{tokens: tokens, projects: projects, hub: hub}
}

func (h *WSHandler) Serve(c *gin.Context) {
	token := c.Query("token")
	projectIDStr := c.Query("projectId")
	if token == "" || projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	claims, err := h.tokens.Parse(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.projects.EnsureMember(c.Request.Context(), projectID, userID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		writeError(c, err)
		return
	}

	conn, err := ws.Upgrade(c.Writer, c.Request)
	if err != nil {
		return
	}

	client := ws.NewClient(conn)
	h.hub.Register(projectIDStr, client)
	defer func() {
		h.hub.Unregister(projectIDStr, client)
		conn.Close()
	}()

	if err := ws.WriteConnected(conn, projectIDStr); err != nil {
		return
	}

	go client.WritePump()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
