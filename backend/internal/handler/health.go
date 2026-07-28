package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health проверка доступности сервера.
//
// @Summary     Проверка health
// @Tags        system
// @Produce     plain
// @Success     200 {string} string "ok"
// @Router      /health [get]
func Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
