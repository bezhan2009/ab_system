package observability

import (
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/http/dto"
	"ab_system/internal/lib/uuid"
	"ab_system/pkg/errs"

	"github.com/gin-gonic/gin"
)

func ValidateQueryAndParamUserID(c *gin.Context) (userID string, err error) {
	userRole := c.GetString(contextkeysHttp.UserRoleCtx)

	if dto.Role(userRole) == dto.RoleAdmin {
		queryUserID := c.Query("userId")
		if queryUserID != "" {
			if err := uuid.ValidateUUID(queryUserID); err != nil {
				return "", errs.ErrIdIsInvalid
			}
			userID = queryUserID
		} else {
			paramUserID := c.Param("id")
			if paramUserID != "" {
				if err := uuid.ValidateUUID(paramUserID); err != nil {
					return "", errs.ErrIdIsInvalid
				}
				userID = paramUserID
			}
		}
	} else {
		userID = c.GetString(contextkeysHttp.UserIDCtx)
	}

	return userID, nil
}
