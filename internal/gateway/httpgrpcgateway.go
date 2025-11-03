package gateway

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"rpc/pkg/api/test"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger *zap.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(lrw, r)

		duration := time.Since(start)

		logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", lrw.statusCode),
			zap.Duration("duration", duration),
			zap.String("user_agent", r.UserAgent()),
		)
	}
}

func StartGateway(ctx context.Context, grpcAddr string, httpAddr string, logger *zap.Logger) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := test.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		mux,
		grpcAddr,
		opts,
	)
	if err != nil {
		logger.Error("Failed to register gateway", zap.Error(err))
		return err
	}

	wrappedMux := wrapLogging(mux, logger)

	logger.Info("Gateway started successfully",
		zap.String("http_addr", httpAddr),
		zap.String("grpc_addr", grpcAddr),
	)

	return http.ListenAndServe(httpAddr, wrappedMux)
}

func wrapLogging(mux *runtime.ServeMux, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingMiddleware(logger, mux.ServeHTTP)(w, r)
	})
}
