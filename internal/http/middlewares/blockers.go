package middlewares

import (
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"

	"github.com/gin-gonic/gin"
)

func BlockViewer(c *gin.Context) {
	userRole := c.GetString(contextkeysHttp.UserRoleCtx)

	if userRole == "viewer" {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}

	c.Next()
}
