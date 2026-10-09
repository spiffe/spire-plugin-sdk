package internal

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unaryPanicRecoveryInterceptor(logger hclog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (_ any, err error) {
		defer recoverCall(logger, info.FullMethod, &err)
		return handler(ctx, req)
	}
}

func streamPanicRecoveryInterceptor(logger hclog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer recoverCall(logger, info.FullMethod, &err)
		return handler(srv, stream)
	}
}

// recoverCall must be deferred on the RPC handler goroutine so that a panic
// fails the call instead of terminating the plugin process.
func recoverCall(logger hclog.Logger, method string, err *error) {
	recovered := recover()
	if recovered == nil {
		return
	}
	// Format the value before passing it to the JSON logger. fmt handles panics
	// from Error/String methods and avoids invoking JSON marshalers.
	logger.Error("Recovered a panic in a plugin call",
		"method", method, "panic", fmt.Sprint(recovered), "stack", string(debug.Stack()))
	*err = status.Error(codes.Internal, "internal error")
}
