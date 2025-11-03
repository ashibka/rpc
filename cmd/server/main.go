package main

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"rpc/internal/config"
	"rpc/internal/gateway"
	"rpc/internal/interceptor"
	"rpc/internal/server"
	"rpc/pkg/api/test"
	"strconv"
	"sync"
	"syscall"
)

func main() {
	//чтение кнфига
	cfg, err := config.ParseConfig("./config/.env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		return
	}

	//инициализируем зап логгер

	logger, err := func(logLevel string) (*zap.Logger, error) {
		if logLevel == "debug" {
			return zap.NewDevelopment()
		}
		return zap.NewProduction()
	}(cfg.LogLevel)
	defer logger.Sync()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
		return
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcserver := grpc.NewServer(grpc.UnaryInterceptor(interceptor.ZapLog(logger)))
	orderServer := server.NewServer()
	reflection.Register(grpcserver)
	test.RegisterOrderServiceServer(grpcserver, orderServer)

	logger.Info("Starting servers",
		zap.String("grpc_port", strconv.Itoa(cfg.Port)),
		zap.String("http_port", strconv.Itoa(cfg.Port)),
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Port))
		if err != nil {
			logger.Fatal("failed to listen on grpc port", zap.Error(err))
		}
		logger.Info("Trying to start grpc order server")
		if err := grpcserver.Serve(listener); err != nil {
			logger.Info("Failed to serve grpc server", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Info("HTTP gateway starting", zap.String("port", strconv.Itoa(cfg.GwPort)))
		if err := gateway.StartGateway(ctx, ":"+strconv.Itoa(cfg.Port), ":"+strconv.Itoa(cfg.GwPort), logger); err != nil {
			logger.Info("Failed to start gateway", zap.Error(err))
		}
	}()
	shdCh := make(chan os.Signal, 1)
	signal.Notify(shdCh, syscall.SIGINT, syscall.SIGTERM)

	<-shdCh
	logger.Info("Gf shutdown started")

	cancel()

	grpcStopped := make(chan struct{})
	go func() {
		grpcserver.GracefulStop()
		close(grpcStopped)
	}()

	select {
	case <-grpcStopped:
		logger.Info("gRPC server stopped")
	case <-ctx.Done():
		logger.Warn("gRPC forced shutdown")
		grpcserver.Stop()
	}

	wg.Wait()
	logger.Info("grpcServer and GW stopped gracefully")
}
