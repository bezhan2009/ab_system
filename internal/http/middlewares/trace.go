package middlewares

import (
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/lib/uuid"
	"ab_system/pkg/errs"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware(c *gin.Context) {
	traceID := c.GetHeader(contextkeysHttp.TraceHeader)
	if traceID == "" {
		var err error
		traceID, err = uuid.NewUUID4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errs.ErrCantGenerateTrace,
			})
			return
		}
	}

	ctx := context.WithValue(c.Request.Context(), contextkeysHttp.TraceIdCtx, traceID)

	c.Request = c.Request.WithContext(ctx)

	c.Set(contextkeysHttp.TraceIdCtx, traceID)

	c.Header(contextkeysHttp.TraceHeader, traceID)

	c.Next()
}
