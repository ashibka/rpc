package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"rpc/internal/config"
	"rpc/internal/interceptor"
	"rpc/internal/server"
	"rpc/pkg/api/test"
	"strconv"
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

	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
		return
	}

	grpcserver := grpc.NewServer(grpc.UnaryInterceptor(interceptor.ZapLog(logger)))
	orderServer := server.NewServer()

	reflection.Register(grpcserver)

	test.RegisterOrderServiceServer(grpcserver, orderServer)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen on %v:%v\n", listener.Addr(), err)
	}
	log.Printf("Trying to start grpc order server on %v\n", listener.Addr())
	if err := grpcserver.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
