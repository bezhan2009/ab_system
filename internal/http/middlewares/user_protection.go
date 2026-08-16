package middlewares

import (
	"ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"

	"github.com/gin-gonic/gin"
)

func UserProtection(c *gin.Context) {
	userId := c.Param("id")

	userIdToken := c.GetString(http.UserIDCtx)
	userRole := c.GetString(http.UserRoleCtx)

	if userRole != "admin" && userIdToken != userId {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}

	c.Next()
}
