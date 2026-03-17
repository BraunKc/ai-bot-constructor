package orchestratorinterceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const tokenCtxKey = "authorization"

func UnaryAuthInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		tokenRaw := ctx.Value(tokenCtxKey)

		if tokenStr, ok := tokenRaw.(string); ok && tokenStr != "" {
			md := metadata.Pairs("authorization", "Bearer "+strings.TrimSpace(tokenStr))
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
