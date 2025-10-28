package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ZapLog(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		startTime := time.Now()
		logger.Info("requested",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			logger.Error("request failed",
				zap.String("method", info.FullMethod),
				zap.Error(err), // логируем ошибку
				zap.Duration("duration", duration),
			)
		} else {
			logger.Info("request completed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Any("response", resp),
			)
		}

		return resp, err
	}
}
