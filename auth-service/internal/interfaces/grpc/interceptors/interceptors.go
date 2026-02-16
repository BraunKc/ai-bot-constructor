package grpcinterceptors

import (
	"context"
	"strings"

	appcontext "github.com/braunkc/ai-bot-constructor/auth-service/internal/application/context"
	"github.com/braunkc/ai-bot-constructor/auth-service/internal/infra/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	tokenManager     *jwt.TokenManager
	protectedMethods map[string]struct{}
}

func NewAuthInterceptor(tokenManager *jwt.TokenManager, protectedMethods map[string]struct{}) *AuthInterceptor {
	return &AuthInterceptor{
		tokenManager:     tokenManager,
		protectedMethods: protectedMethods,
	}
}

func (ai *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if _, protected := ai.protectedMethods[info.FullMethod]; !protected {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization headers")
		}

		authHeader := authHeaders[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := ai.tokenManager.VerifyToken(accessToken)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		return handler(appcontext.ContextWithUserID(ctx, userID), req)
	}
}
