package middlewares

import (
	"ab_system/internal/http/response"
	"ab_system/internal/lib/uuid"
	"ab_system/pkg/errs"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateUUID(c *gin.Context) {
	for _, param := range c.Params {
		paramLower := strings.ToLower(param.Key)
		if strings.Contains(paramLower, "id") && paramLower != "" {
			if err := uuid.ValidateUUID(param.Value); err != nil {
				response.HandleError(c, errs.ErrIdIsInvalid)

				return
			}
		}
	}

	c.Next()
}
