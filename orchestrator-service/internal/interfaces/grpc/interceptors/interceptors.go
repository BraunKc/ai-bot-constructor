package grpcinterceptors

import (
	"context"
	"strings"

	"github.com/braunkc/ai-bot-constructor/orchestrator-service/internal/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(noRequire map[string]struct{}, tokenManager *auth.TokenManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := noRequire[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		accessToken := strings.TrimPrefix(authHeaders[0], "Bearer ")

		userID, err := tokenManager.VerifyToken(accessToken)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "user_id", userID)

		return handler(ctx, req)
	}
}
