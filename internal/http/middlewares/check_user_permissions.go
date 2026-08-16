package middlewares

import (
	"ab_system/internal/domain/models"
	"ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"

	"github.com/gin-gonic/gin"
)

func CheckUserAdmin(c *gin.Context) {
	userRole := c.GetString(http.UserRoleCtx)

	if models.Role(userRole) != models.RoleAdmin {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}

	c.Set(http.UserAdminCtx, 1)
	c.Next()
}

func CheckUserAdminSoft(c *gin.Context) {
	userRole := c.GetString(http.UserRoleCtx)
	if models.Role(userRole) == models.RoleAdmin {
		c.Set(http.UserAdminCtx, true)
	} else {
		c.Set(http.UserAdminCtx, false)
	}
	c.Next()
}

func CheckUserExperimenter(c *gin.Context) {
	if c.GetBool(http.UserAdminCtx) {
		c.Next()
		return
	}
	userRole := c.GetString(http.UserRoleCtx)
	if models.Role(userRole) != models.RoleExperimenter {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}
	c.Next()
}

func CheckUserApprover(c *gin.Context) {
	if c.GetBool(http.UserAdminCtx) {
		c.Next()
		return
	}
	userRole := c.GetString(http.UserRoleCtx)
	if models.Role(userRole) != models.RoleApprover {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}
	c.Next()
}

func CheckUserViewer(c *gin.Context) {
	if c.GetBool(http.UserAdminCtx) {
		c.Next()
		return
	}
	userRole := c.GetString(http.UserRoleCtx)
	if models.Role(userRole) != models.RoleViewer {
		response.HandleError(c, errs.ErrPermissionDenied)
		return
	}
	c.Next()
}