package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/planeweave/planeweave/internal/domain"
)

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest), errors.Is(err, domain.ErrPasswordTooShort):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errorCode(err)})
	case errors.Is(err, domain.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorCode(err)})
	case errors.Is(err, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrInviteNotFound), errors.Is(err, domain.ErrInviteExpired):
		c.JSON(http.StatusNotFound, gin.H{"error": errorCode(err)})
	case errors.Is(err, domain.ErrCycleDetected):
		c.JSON(http.StatusConflict, gin.H{"error": "cycle_detected"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
	}
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		return "invalid_request"
	case errors.Is(err, domain.ErrPasswordTooShort):
		return "password_too_short"
	case errors.Is(err, domain.ErrInvalidCredentials):
		return "invalid_credentials"
	case errors.Is(err, domain.ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, domain.ErrInviteNotFound):
		return "invite_not_found"
	case errors.Is(err, domain.ErrInviteExpired):
		return "invite_expired"
	default:
		return "not_found"
	}
}

func parseUUID(c *gin.Context, param string) (string, bool) {
	value := c.Param(param)
	if _, err := domain.StringToUUID(value); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid_request"})
		return "", false
	}
	return value, true
}
