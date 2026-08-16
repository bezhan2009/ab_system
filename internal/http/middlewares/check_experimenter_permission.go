package middlewares

import (
	"ab_system/internal/domain/repository"
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/http/response"
	"ab_system/pkg/errs"

	"github.com/gin-gonic/gin"
)

func CheckExperimenterPermission(experimentReader repository.ExperimentReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(contextkeysHttp.UserIDCtx)
		expID := c.Param("id")

		experiment, err := experimentReader.GetExperimentByIdArchive(c.Request.Context(), expID)
		if err != nil {
			response.HandleError(c, err)

			return
		}

		if experiment.OwnerID != userID {
			response.HandleError(c, errs.ErrPermissionDenied)

			return
		}

		c.Next()
	}
}
