package interceptors

import (
	contextkeysHttp "ab_system/internal/http"
	"ab_system/internal/lib/uuid"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TraceUnaryInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	traceID := ""
	if vals := md.Get(contextkeysHttp.TraceHeader); len(vals) > 0 {
		traceID = vals[0]
	}
	if traceID == "" {
		var err error
		traceID, err = uuid.NewUUID4()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cant generate trace id")
		}
	}

	ctx = context.WithValue(ctx, contextkeysHttp.TraceIdCtx, traceID)

	_ = grpc.SetHeader(ctx, metadata.Pairs(contextkeysHttp.TraceHeader, traceID))

	return handler(ctx, req)
}
