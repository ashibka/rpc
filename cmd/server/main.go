package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	redislib "github.com/redis/go-redis/v9"
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
	"rpc/internal/repository/cached"
	"rpc/internal/repository/postgres"
	redisrepo "rpc/internal/repository/redis"
	"rpc/internal/server"
	"rpc/pkg/api/test"
	"strconv"
	"sync"
	"syscall"
	"time"
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

	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Ping(ctxTimeout); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	fmt.Println("Postgres connected sucssefully")

	redisClient := redislib.NewClient(&redislib.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	defer redisClient.Close()

	if err := redisClient.Ping(ctxTimeout).Err(); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	orderRepo := postgres.NewOrderRepository(db)
	redisRepo := redisrepo.NewOrderRepository(redisClient)
	cachedRepo := cached.NewCachedRepository(redisRepo, orderRepo)

	grpcserver := grpc.NewServer(grpc.UnaryInterceptor(interceptor.ZapLog(logger)))
	orderServer := server.NewServer(cachedRepo)
	reflection.Register(grpcserver)
	test.RegisterOrderServiceServer(grpcserver, orderServer)

	logger.Info("Starting servers",
		zap.String("grpc_port", strconv.Itoa(cfg.Port)),
		zap.String("http_port", strconv.Itoa(cfg.Port)),
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(cfg.Port))
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
		if err := gateway.StartGateway(ctx, "0.0.0.0:"+strconv.Itoa(cfg.Port), "0.0.0.0:"+strconv.Itoa(cfg.GwPort), logger); err != nil {
			logger.Info("Failed to start gateway", zap.Error(err))
		}
	}()
	shdCh := make(chan os.Signal, 1)
	signal.Notify(shdCh, syscall.SIGINT, syscall.SIGTERM)

	<-shdCh
	logger.Info("Gf shutdown started")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	cancel()

	grpcStopped := make(chan struct{})
	go func() {
		grpcserver.GracefulStop()
		close(grpcStopped)
	}()

	select {
	case <-grpcStopped:
		logger.Info("gRPC server stopped")
	case <-shutdownCtx.Done():
		logger.Warn("gRPC forced shutdown: timeout exceeded")
		grpcserver.Stop()
	}

	wg.Wait()
	logger.Info("grpcServer and GW stopped gracefully")
}
