package middlewares

import (
	"ab_system/internal/domain/repository"
	"ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"
	"errors"

	"github.com/gin-gonic/gin"
)

func UserActiveMiddleware(checker repository.UserStatusChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(http.UserIDCtx)

		if userID == "" {
			c.Next()
			return
		}

		active, err := checker.IsActive(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, errs.ErrRecordNotFound) {
				response.HandleError(c, errs.ErrUserNotFound)
				return
			}

			response.HandleError(c, err)
			return
		}

		if !active {
			response.HandleError(c, errs.ErrPermissionDenied)
			return
		}

		c.Next()
	}
}
