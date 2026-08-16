package observability

import (
	"ab_system/internal/http"
	"context"
)

func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(http.TraceIdCtx).(string); ok {
		return traceID
	}

	return ""
}
